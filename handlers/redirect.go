package handlers

import (
	"net/http"
	"strings"
	"url-shortener/storage"
)

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/")
	if shortID == "" {
		http.Error(w, "short url not found", http.StatusNotFound)
		return
	}

	// Проверяем, есть ли такой shortID в хранилище
	longURL, exists := storage.GetURL(shortID)
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	// Делаем редирект на оригинальный URL
	http.Redirect(w, r, longURL, http.StatusFound)
}
