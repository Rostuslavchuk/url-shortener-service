package save

import (
	"errors"
	"log/slog"
	"net/http"

	resp "url_shortener/internal/lib/api/response"
	"url_shortener/internal/lib/random"
	"url_shortener/internal/lib/sl"
	"url_shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLen = 6

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const op = "handlers.save.new"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(req.Context())),
		)

		var decodedReq Request
		err := render.DecodeJSON(req.Body, &decodedReq)
		if err != nil {
			log.Error("faild to decode request body", sl.Err(err))
			render.JSON(w, req, resp.Error("faild to decode request body"))
			return
		}

		if err := validator.New().Struct(decodedReq); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			render.JSON(w, req, resp.Validate(validateErr))
			return
		}

		alias := decodedReq.Alias
		if alias == "" {
			for {
				aliasNew := random.NewRandomString(aliasLen)
				_, err := urlSaver.GetURL(alias)

				if errors.Is(err, storage.ErrNotFound) {
					alias = aliasNew
					break
				}
				if err == nil {
					continue
				}
			}
		}

		_, err = urlSaver.SaveURL(decodedReq.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", decodedReq.URL))
				render.JSON(w, req, resp.Error("url already exists"))
				return
			}
			log.Error("faild to add url", sl.Err(err))
			render.JSON(w, req, resp.Error("faild to add url"))
			return
		}

		render.JSON(w, req, Response{
			Alias:    alias,
			Response: *resp.OK(),
		})
	}
}
