package storage

var (
	storage = make(map[string]string)
)

func SaveURL(shortID, longURL string) {
	storage[shortID] = longURL
}

func GetURL(shortID string) (string, bool) {
	longURL, exists := storage[shortID]
	return longURL, exists
}
