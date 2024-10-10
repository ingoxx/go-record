package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"

	admissv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"log"
	"net/http"
)

// admitFunc 定义了处理Webhook请求的函数类型
type admitFunc func(admissionReview *admissv1.AdmissionReview) *admissv1.AdmissionResponse

// admitForDebug 是一个示例admitFunc，仅打印接收到的请求并允许所有操作
func admitWithReplicaLimit(ar *admissv1.AdmissionReview) *admissv1.AdmissionResponse {
	resourceReq := ar.Request.Resource.Resource
	schemeRos := admissv1.SchemeGroupVersion.WithResource("deployments").Resource

	if schemeRos != resourceReq {
		return &admissv1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Message: "Not a Deployment resource.",
			},
		}
	}

	decoder := serializer.NewSerializerWithOptions(serializer.DefaultMetaFactory, runtime.NewScheme(), runtime.NewScheme(), serializer.SerializerOptions{
		Pretty: true,
	})
	obj, _, err := decoder.Decode(ar.Request.Object.Raw, nil, &appsv1.Deployment{})
	if err != nil {
		return &admissv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: fmt.Sprintf("Could not deserialize request object: %v", err),
			},
		}
	}

	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return &admissv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: "Deserialized object is not a Deployment.",
			},
		}
	}

	maxReplicas := int32(3)
	if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas > maxReplicas {
		message := fmt.Sprintf("The number of replicas (%d) exceeds the maximum allowed (%d).", *deployment.Spec.Replicas, maxReplicas)
		return &admissv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: message,
			},
		}
	}

	return &admissv1.AdmissionResponse{
		Allowed: true,
	}
}

// toAdmissionReview 解析HTTP请求体为AdmissionReview对象
func toAdmissionReview(r io.ReadCloser) (*admissv1.AdmissionReview, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("can't read body: %v", err)
	}
	ar := admissv1.AdmissionReview{}
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("can't unmarshal body: %v", err)
	}
	return &ar, nil
}

// serveHTTP 处理HTTP请求，调用admitFunc处理Webhook
func serveHTTP(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var reviewResponse *admissv1.AdmissionReview
	if r.URL.Path != "/validate" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method type.", http.StatusMethodNotAllowed)
		return
	}

	ar, err := toAdmissionReview(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing AdmissionReview: %v", err), http.StatusBadRequest)
		return
	}

	reviewResponse = &admissv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/admissv1",
		},
		Response: &admissv1.AdmissionResponse{
			UID:     ar.Request.UID,
			Allowed: false, // 默认拒绝，admitFunc将根据实际情况修改
		},
	}

	if ar.Request != nil {
		reviewResponse.Response = admit(ar)
	}

	resp, err := json.Marshal(reviewResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't encode response: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(resp); err != nil {
		http.Error(w, fmt.Sprintf("Can't write response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12, // 设置最低TLS版本
	}

	// 加载证书和私钥
	cert, err := tls.LoadX509KeyPair("", "")
	if err != nil {
		log.Fatalf("Failed to load certificate: %s", err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		serveHTTP(w, r, admitWithReplicaLimit)
	})

	// 启动HTTPS服务器
	srv := &http.Server{
		Addr:      ":8443",   // HTTPS默认端口为443，这里使用8443仅为示例
		TLSConfig: tlsConfig, // 使用配置好的TLS
		Handler:   mux,
	}

	fmt.Println("Listening on :8443...")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
