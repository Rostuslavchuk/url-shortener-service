package delete

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url_shortener/internal/lib/api/response"
	"url_shortener/internal/lib/sl"
	"url_shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) (int64, error)
}

type Request struct {
	Alias string `json:"alias" validate:"required"`
}
type Response struct {
	resp.Response
	Message string `json:"message" validate:"required"`
}

func New(log *slog.Logger, urldeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const op = "handlers.delete.New"

		log = slog.With(
			"op", op,
			"request_id", middleware.GetReqID(req.Context()),
		)

		decodedReq := &Request{
			Alias: chi.URLParam(req, "alias"),
		}
		if decodedReq.Alias == "" {
			log.Info("alias is empty")
			render.JSON(w, req, resp.Error("alias is empty"))
			return
		}

		_, err := urldeleter.DeleteURL(decodedReq.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Error("not found", sl.Err(err))
				render.JSON(w, req, resp.Error("not found"))
				return
			}
			log.Error("faild to delete", sl.Err(err))
			render.JSON(w, req, resp.Error("faild to delete"))
			return
		}

		render.JSON(w, req, &Response{
			Response: *resp.OK(),
			Message:  "deleted successfully",
		})
	}
}
