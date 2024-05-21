package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type PayloadJSON struct {
	Text   string
	Number int
}

func main() {
	var opt int
	for {
		fmt.Println("1. Get")
		fmt.Println("2. Post")
		fmt.Println("3. Exit")
		fmt.Scanf("%d\n", &opt)

		if opt == 1 {
			getData()
		} else if opt == 2 {
			postData()
		} else if opt == 3 {
			break
		} else {
			println("Invalid!")
		}
	}
}

func getData() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		fmt.Println("Cant make request")
		return
	}
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Server is not responsive")
		} else {
			fmt.Println(err)
		}
		return
	}

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}
		fmt.Println(string(body))

	}

	_ = resp.Body.Close()
}

func postData() {
	var opt int
	for {
		fmt.Println("1. Post JSON")
		fmt.Println("2. Post multipart form")
		fmt.Println("3. Cancel")
		fmt.Scanf("%d\n", &opt)
		if opt == 1 {
			var text string
			var num int
			fmt.Print("Enter text: ")
			fmt.Scanf("%s\n", &text)
			fmt.Print("Enter num: ")
			fmt.Scanf("%d\n", &num)

			buf := new(bytes.Buffer)
			pJ := PayloadJSON{Text: text, Number: num}
			err := json.NewEncoder(buf).Encode(&pJ)
			if err != nil {
				fmt.Println("Failed making JSON payload")
			}
			// fmt.Println(string(buf.Bytes()))

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080", buf)
			if err != nil {
				fmt.Println("Cant make request")
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
			}

			defer resp.Body.Close()
			return

		} else if opt == 2 {
			var desc string
			fmt.Print("Enter description: ")
			fmt.Scanf("%s\n", &desc)
			reqBody := new(bytes.Buffer)
			w := multipart.NewWriter(reqBody)

			for k, v := range map[string]string{
				"date":        time.Now().Format(time.RFC3339),
				"description": desc,
			} {
				err := w.WriteField(k, v)
				if err != nil {
					fmt.Println(err)
				}
			}

			for i, file := range []string{
				"./files/hello.txt",
				"./files/test.txt",
			} {
				filePart, err := w.CreateFormFile(fmt.Sprintf("file%d", i+1),
					filepath.Base(file))
				if err != nil {
					fmt.Println(err)
				}
				f, err := os.Open(file)
				if err != nil {
					fmt.Println(err)
				}
				_, err = io.Copy(filePart, f)
				_ = f.Close()
				if err != nil {
					fmt.Println(err)
				}
			}
			err := w.Close()
			if err != nil {
				fmt.Println(err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:8080", reqBody)
			if err != nil {
				fmt.Println("Cant make request")
				return
			}
			req.Header.Set("Content-Type", w.FormDataContentType())

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
			}

			defer resp.Body.Close()
			return

		} else if opt == 3 {
			return
		} else {
			fmt.Println("Invalid!")
		}

	}
}
