//go:build !solution

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type URLStorage struct {
	mutex sync.Mutex
	data  map[string]string
}

type URLShortener struct {
	storage *URLStorage
}

type HTTPHandler struct {
	shortener *URLShortener
}

func NewURLStorage() *URLStorage {
	return &URLStorage{data: make(map[string]string)}
}

func (storage *URLStorage) SaveKey(key, url string) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	storage.data[key] = url
}

func (storage *URLStorage) RetrieveURL(key string) (string, bool) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	url, exists := storage.data[key]
	return url, exists
}

func NewURLShortener(storage *URLStorage) *URLShortener {
	return &URLShortener{storage: storage}
}

func (shortener *URLShortener) GenerateKey(url string) (string, error) {
	hasher := sha1.New()
	if _, err := hasher.Write([]byte(url)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil))[:10], nil
}

func (shortener *URLShortener) ShortenURL(url string) (string, error) {
	key, err := shortener.GenerateKey(url)
	if err != nil {
		return "", err
	}
	shortener.storage.SaveKey(key, url)
	return key, nil
}

func NewHTTPHandler(shortener *URLShortener) *HTTPHandler {
	return &HTTPHandler{shortener: shortener}
}

func (handler *HTTPHandler) handleShortenRequest(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		respondWithError(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	input, err := parseJSONRequest(request)
	if err != nil || input["url"] == "" {
		respondWithError(writer, "Invalid Request", http.StatusBadRequest)
		return
	}

	url := input["url"]
	key, err := handler.shortener.ShortenURL(url)
	if err != nil {
		respondWithError(writer, "Server Error", http.StatusInternalServerError)
		return
	}

	output := map[string]string{"url": url, "key": key}
	respondWithJSON(writer, output, http.StatusOK)
}

func (handler *HTTPHandler) handleRedirectRequest(writer http.ResponseWriter, request *http.Request) {
	key := extractKeyFromPath(request.URL.Path, "/go/")
	if key == "" {
		respondWithError(writer, "Key Not Found", http.StatusNotFound)
		return
	}

	url, exists := handler.shortener.storage.RetrieveURL(key)
	if !exists {
		respondWithError(writer, "Key Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(writer, request, url, http.StatusFound)
}

func parseJSONRequest(request *http.Request) (map[string]string, error) {
	var data map[string]string
	err := json.NewDecoder(request.Body).Decode(&data)
	return data, err
}

func extractKeyFromPath(path, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}

func respondWithJSON(writer http.ResponseWriter, data interface{}, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		return
	}
}

func respondWithError(writer http.ResponseWriter, message string, statusCode int) {
	http.Error(writer, message, statusCode)
}

func main() {
	port := flag.String("port", "8080", "HTTP server port")
	flag.Parse()

	storage := NewURLStorage()
	shortener := NewURLShortener(storage)
	handler := NewHTTPHandler(shortener)

	http.HandleFunc("/shorten", handler.handleShortenRequest)
	http.HandleFunc("/go/", handler.handleRedirectRequest)

	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Server is running on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
