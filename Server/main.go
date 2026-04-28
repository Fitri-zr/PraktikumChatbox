package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// request dari client
type Request struct {
	Text string `json:"text"`
}

// response ke client
type Response struct {
	Reply string `json:"reply"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	// ✅ CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// preflight
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	// decode request
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Text == "" {
		http.Error(w, "Request tidak valid", http.StatusBadRequest)
		return
	}

	// ambil API key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		http.Error(w, "API KEY belum di-set", http.StatusInternalServerError)
		return
	}

	// URL Gemini
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash-lite:generateContent?key=" + apiKey

	// prompt khusus IoT
	prompt := "Anda adalah chatbot khusus IoT. Jawab hanya tentang IoT.\n\nPertanyaan: " + req.Text

	// body request
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(body)

	// request ke Gemini
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Gagal konek ke Gemini", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	// 🔍 DEBUG (lihat di docker logs)
	fmt.Println("===== GEMINI RESPONSE =====")
	fmt.Println(string(responseBody))
	fmt.Println("===========================")

	// cek status
	if resp.StatusCode != 200 {
		http.Error(w, "API Gemini error", http.StatusInternalServerError)
		return
	}

	// parsing JSON
	var gemini map[string]interface{}
	json.Unmarshal(responseBody, &gemini)

	reply := "Maaf, terjadi kesalahan."

	// parsing aman
	if candidates, ok := gemini["candidates"].([]interface{}); ok {
		for _, c := range candidates {
			candidate, ok := c.(map[string]interface{})
			if !ok {
				continue
			}

			content, ok := candidate["content"].(map[string]interface{})
			if !ok {
				continue
			}

			parts, ok := content["parts"].([]interface{})
			if !ok {
				continue
			}

			for _, p := range parts {
				part, ok := p.(map[string]interface{})
				if !ok {
					continue
				}

				if text, ok := part["text"].(string); ok && text != "" {
					reply = text
					break
				}
			}
		}
	}

	// kirim ke frontend
	json.NewEncoder(w).Encode(Response{
		Reply: reply,
	})
}

func main() {
	http.HandleFunc("/chat", chatHandler)

	fmt.Println("Server jalan di http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}