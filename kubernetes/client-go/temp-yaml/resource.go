package main

import (
	"bytes"
	"fmt"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"text/template"
)

type AppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of App. Edit app_types.go to remove/update
	EnableIngress bool   `json:"enable_ingress,omitempty"`
	EnableService bool   `json:"enable_service"`
	Replicas      int32  `json:"replicas"`
	Image         string `json:"image"`
}

type App struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AppSpec `json:"spec,omitempty"`
}

func parseYamlString(templateName string, app *App) ([]byte, error) {
	deployment := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ObjectMeta.Name}}-deployment
  namespace: {{.ObjectMeta.Namespace}}
  labels:
    app: {{.ObjectMeta.Name}}-app
spec:
  replicas: {{.Spec.Replicas}}
  selector:
    matchLabels:
      app: {{.ObjectMeta.Name}}-app
  template:
    metadata:
      labels:
        app: {{.ObjectMeta.Name}}-app
    spec:
      containers:
        - name: {{.ObjectMeta.Name}}-pod
          image: {{.Spec.Image}}
          ports:
            - containerPort: 80`

	service := `apiVersion: v1
				kind: Service
				metadata:
				  name: {{.ObjectMeta.Name}}-svc
				  namespace: {{.ObjectMeta.Namespace}}
				spec:
				  ports:
					- port: 8081
					  targetPort: 80
					  protocol: TCP
				  type: NodePort
				  selector:
					app: {{.ObjectMeta.Name}}-app`

	var yamlContent = map[string]string{
		"deployment": deployment,
		"service":    service,
	}

	tmpl, err := template.New(templateName).Parse(yamlContent[templateName])
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, app); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}

func parseYaml(templateName string, app *App) ([]byte, error) {
	dir, _ := os.Getwd()
	fmt.Println("dir >>>> ", dir)
	file, _ := os.ReadFile(fmt.Sprintf("./template/%s.yaml", templateName))
	tmpl, err := template.New(templateName).Parse(string(file))
	if err != nil {
		return []byte{}, err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, app); err != nil {
		return []byte{}, err
	}

	fmt.Println(" >>>>>>>> ", string(buf.Bytes()))

	return buf.Bytes(), nil
}

func NewDeployment(app *App) (*appsV1.Deployment, error) {
	var dep = new(appsV1.Deployment)
	b, err := parseYamlString("deployment", app)
	if err != nil {
		return nil, err
	}

	fmt.Println("tmlp >>> ", string(b))

	err = yaml.Unmarshal(b, dep)
	if err != nil {
		return nil, err
	}

	fmt.Println("dep >>> ", dep)

	return dep, nil
}

func NewService(app *App) (*coreV1.Service, error) {
	var svc = &coreV1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name + "-" + "svc",
			Namespace: app.Namespace,
		},
		Spec: coreV1.ServiceSpec{
			Selector: map[string]string{
				"app": app.Name + "-" + "app",
			},
			Type: coreV1.ServiceType("NodePort"),
			Ports: []coreV1.ServicePort{
				{
					Name:     "http",
					Protocol: coreV1.Protocol("TCP"),
					Port:     8081,
					NodePort: 30006,
					TargetPort: intstr.IntOrString{
						IntVal: 80,
					},
				},
			},
		},
	}

	//b, err := parseYaml("service", app)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err := yaml.UnmarshalStrict(b, svc); err != nil {
	//	return nil, err
	//}

	return svc, nil
}

func main() {
	appSpec := AppSpec{
		Image:    "nginx",
		Replicas: 2,
	}
	app := &App{
		Spec: appSpec,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "web",
		},
	}

	deployment, err := NewDeployment(app)
	fmt.Println(deployment, err)
}
