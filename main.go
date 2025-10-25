package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

type Response struct {
	Secrets map[string]string `json:"secrets"`
}

var secretDir = "/data/secrets"

func main() {
	if env := os.Getenv("SECRETS_DIR"); env != "" {
		secretDir = env
	}

	r := mux.NewRouter()
	r.HandleFunc("/webhook", webhookHandler).Methods(http.MethodGet)

	addr := ":8080"
	fmt.Printf("ESO FileSecret Provider running on %s (serving from %s)\n", addr, secretDir)
	log.Fatal(http.ListenAndServe(addr, r))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing 'key' query parameter", http.StatusBadRequest)
		return
	}

	// TODO: Verify key is safe string

	keyFile := fmt.Sprintf("%s.json", key)
	filePath := filepath.Join(secretDir, keyFile)

	content, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Secret not found: %s", key), http.StatusNotFound)
		return
	}

	var secrets map[string]string
	if err := json.Unmarshal(content, &secrets); err != nil {
		http.Error(w, "Invalid secrets data in secret file", http.StatusInternalServerError)
		return
	}

	response := Response{Secrets: secrets}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		log.Printf("error encoding JSON response: %v", err)
		return
	}
}
