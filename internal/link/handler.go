package link

import (
	"fmt"
	"net/http"
	"strconv"
	"url_shortener/configs"
	"url_shortener/pkg/event"
	"url_shortener/pkg/middleware"
	"url_shortener/pkg/req"
	"url_shortener/pkg/res"

	"gorm.io/gorm"
)

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}
type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.Handle("POST /link", middleware.IsAuthed(handler.CreateLink(), deps.Config))
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.UpdLink(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.DeleteLink(), deps.Config))
	router.Handle("GET /link", middleware.IsAuthed(handler.GetAll(), deps.Config))

}

func (handler *LinkHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.GetByID(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.MakeJson(w, nil, 200)
	}
}

func (handler *LinkHandler) UpdLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		e, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if ok {
			fmt.Println(e)
		}

		body, err := req.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			return
		}

		idString := r.PathValue("id")

		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.MakeJson(w, link, 201)

	}
}

func (handler *LinkHandler) CreateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LinkCreateRequest](&w, r)
		if err != nil {
			return
		}

		link := NewLink(body.Url)
		for {
			existedLink, _ := handler.LinkRepository.GetByHash(link.Hash)
			if existedLink.Hash == "" {
				break
			}
			link.GenerateHash()
		}

		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.MakeJson(w, createdLink, 201)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		link, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// handler.StatRepository.AddClick(link.ID)
		go handler.EventBus.Publish(event.Event{
			Type: event.LinkVisitedEvent,
			Data: link.ID,
		})
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.Count()
		res.MakeJson(w, GetAllLinksResponse{Links: links, Count: int(count)}, 200)
	}
}
