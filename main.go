package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Message string `json:"message"`
	XResult string `json:"x-result"`
	XBody   string `json:"x-body"`
}

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "x-test,ngrok-skip-browser-warning,Content-Type,Accept,Access-Control-Allow-Headers")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	resp := Response{
		Message: "269018",
		XResult: r.Header.Get("x-test"),
		XBody:   string(body),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/result4/", enableCORS(resultHandler))
	http.HandleFunc("/healthz", healthzHandler) // или на "/"

	log.Printf("Listening on :%s…", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
