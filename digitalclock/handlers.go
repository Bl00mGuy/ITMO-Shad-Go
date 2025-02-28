package main

import (
	"image/png"
	"net/http"
)

func handleClockRequest(writer http.ResponseWriter, request *http.Request) {
	clockRequest, err := ParseClockRequest(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	img := GenerateClockImage(clockRequest)

	writer.Header().Set("Content-Type", "image/png")
	if err := png.Encode(writer, img); err != nil {
		http.Error(writer, "failed to encode image", http.StatusInternalServerError)
	}
}
