package stat

import (
	"net/http"
	"time"
	"url_shortener/configs"
	"url_shortener/pkg/middleware"
	"url_shortener/pkg/res"
)

type StatHandler struct {
	StatRepository *StatRepository
}
type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}

	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStats(), deps.Config))
}

func (handler StatHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "Invalid from", http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to", http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")
		if by != "month" && by != "day" {
			http.Error(w, "Invalid 'by' value. Allowed: 'month' or 'day'", http.StatusBadRequest)
			return
		}
		stats := handler.StatRepository.GetStats(from, to, by)
		res.MakeJson(w, stats, 200)
	}
}
