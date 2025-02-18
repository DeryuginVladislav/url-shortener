package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"url-shortener/storage"
	"url-shortener/util"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var u ShortenRequest
	err = json.Unmarshal(body, &u)
	if err != nil || u.URL == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Генерируем короткий ID
	shortID := util.GenerateShortID()
	shortURL := "http://localhost:8000/" + shortID

	//сохраняем
	storage.SaveURL(shortID, u.URL)

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ShortenResponse{ShortURL: shortURL})
}
