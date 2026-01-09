package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"url_shortener/internal/lib/sl"
	"url_shortener/internal/storage"

	resp "url_shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Alias string `json:"alias" validate:"required"`
}
type Response struct {
	resp.Response
	URL string `json:"url" validate:"required"`
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const op = "handlers.redirect.new"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(req.Context())),
		)

		decodedReq := &Request{
			Alias: chi.URLParam(req, "alias"),
		}

		if decodedReq.Alias == "" {
			log.Info("alias is empty")
			render.JSON(w, req, resp.Error("invalid request"))
			return
		}

		url, err := urlGetter.GetURL(decodedReq.Alias)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Error("not found", sl.Err(err), slog.String("alias", decodedReq.Alias))
				render.JSON(w, req, resp.Error("not found"))
				return
			}
			log.Error("faild to get url", sl.Err(err))
			render.JSON(w, req, resp.Error("faild to get url"))
			return
		}

		http.Redirect(w, req, url, http.StatusFound)
	}
}
