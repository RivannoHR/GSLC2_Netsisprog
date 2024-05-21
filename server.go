package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PayloadJSON struct {
	Text   string
	Number int
}

func main() {
	fmt.Println("Server listening on port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprintf(w, "i am response body")
		} else if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
		} else if r.Method == http.MethodPost {
			contentType := r.Header.Get("Content-Type")
			switch {
			case strings.HasPrefix(contentType, "application/json"):
				handleJSON(w, r)
			case strings.HasPrefix(contentType, "multipart/form-data"):
				handleMultipartForm(w, r)
			}
		}
	})
	http.ListenAndServe(":8080", nil)
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	var pJ PayloadJSON
	err := json.NewDecoder(r.Body).Decode(&pJ)
	if err != nil {
		fmt.Println("Decode Failed")
		return
	}
	fmt.Println("JSON received:", pJ)
}

func handleMultipartForm(w http.ResponseWriter, r *http.Request) {
	reader, err := r.MultipartReader()
	if err != nil {
		fmt.Println("Error parsing multipart form:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading multipart form part:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if part.FormName() != "" {
			value, err := io.ReadAll(part)
			if err != nil {
				fmt.Println("Error reading form field data:", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Printf("Form field: %s - Value: %s\n", part.FormName(), string(value))
			// Handle the form field value based on its name (e.g., "date", "description")
		}
	}
}
