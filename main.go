package main

import (
	"encoding/json"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/go-chi/chi/v5"
)

type Page struct {
	Title string
}

type URL struct {
	Url string `json:"url"`
}

func main() {
	r := chi.NewRouter()
	r.Post("/generator", viewCodeHandler)

	http.ListenAndServe(":8080", r)
}

func viewCodeHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("impossible to read body request: %s", err)
	}

	var dataString URL
	if err = json.Unmarshal(requestBody, &dataString); err != nil {
		log.Fatalf("impossible to convert body request: %s", err)
	}

	qrCode, err := qr.Encode(dataString.Url, qr.L, qr.Auto)
	if err != nil {
		log.Fatalf("impossible to Encode QR: %s", err)
	}

	qrCode, err = barcode.Scale(qrCode, 512, 512)
	if err != nil {
		log.Fatalf("impossible to Scale QR: %s", err)
	}

	file, err := os.Create("barcode.png")
	if err != nil {
		log.Fatalf("impossible to create file: %s", err)
	}
	defer file.Close()

	err = png.Encode(file, qrCode)
	if err != nil {
		log.Fatalf("impossible to encode barcode in PNG: %s", err)
	}

	fileBytes, err := ioutil.ReadFile(file.Name())
	if err != nil {
		log.Fatalf("impossible to read PNG: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}
