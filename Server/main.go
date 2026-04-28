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

// response ke client (biar simpel)
type Response struct {
	Reply string `json:"reply"`
}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	// ✅ CORS (WAJIB untuk browser)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// ✅ Handle preflight
	if r.Method == "OPTIONS" {
		return
	}

	// ✅ hanya POST
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

	// ambil API KEY
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		http.Error(w, "API KEY belum di-set", http.StatusInternalServerError)
		return
	}

	// URL Gemini
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash-lite:generateContent?key=" + apiKey

	// ✅ PROMPT KHUSUS IOT
	prompt := "Anda adalah chatbot khusus IoT (Internet of Things). Jawab hanya pertanyaan terkait IoT. Jika di luar IoT, jawab: 'Maaf, saya hanya membahas IoT.'\n\nPertanyaan: " + req.Text

	// body ke Gemini
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

	// kirim ke Gemini
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Gagal konek ke API Gemini", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// baca response
	responseBody, _ := io.ReadAll(resp.Body)

	// parsing JSON Gemini
	var gemini map[string]interface{}
	json.Unmarshal(responseBody, &gemini)

	// ambil teks jawaban
	reply := "Maaf, terjadi kesalahan."

	if candidates, ok := gemini["candidates"].([]interface{}); ok && len(candidates) > 0 {
		candidate := candidates[0].(map[string]interface{})

		if content, ok := candidate["content"].(map[string]interface{}); ok {
			if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
				part := parts[0].(map[string]interface{})
				if text, ok := part["text"].(string); ok {
					reply = text
				}
			}
		}
	}

	// kirim ke client (clean JSON)
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