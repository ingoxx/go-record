package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// GitHubPushEvent represents the structure of a GitHub push event.
type GitHubPushEvent struct {
	Ref    string `json:"ref"`
	Before string `json:"before"`
	After  string `json:"after"`
	Repo   struct {
		Name string `json:"name"`
	} `json:"repository"`
	Pusher struct {
		Name string `json:"name"`
	} `json:"pusher"`
}

// verifySignature verifies the signature of the incoming request.
func verifySignature(r *http.Request, secret string) bool {
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		return false
	}
	expectedMAC := []byte(signature[7:]) // Strip "sha256=" prefix
	mac := hmac.New(sha256.New, []byte(secret))
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return false
	}
	mac.Write(body)
	actualMAC := mac.Sum(nil)
	return hmac.Equal(expectedMAC, actualMAC)
}

// handleWebhook handles the incoming GitHub webhook.
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	secret := "123456" // 替换为你的GitHub Webhook密钥
	if !verifySignature(r, secret) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	reader := bufio.NewReader(r.Body)
	var b = make([]byte, 1024)
	var data = make([]byte, 1024)
	for {
		n, err := reader.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		data = append(data, b[:n]...)
	}

	var event GitHubPushEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Received push event for repository: %s, ref: %s, by: %s\n",
		event.Repo.Name, event.Ref, event.Pusher.Name)
}

func main() {
	http.HandleFunc("/webhook", handleWebhook)
	log.Println("Starting server on :8089")
	if err := http.ListenAndServe(":8089", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

}
