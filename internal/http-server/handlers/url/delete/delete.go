package delete

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type urlDeleter interface {
	DeleteURL(alias string) error
}

// New - создаёт и возвращает готовый обработчик для удаления url.
func New(log *slog.Logger, urlDeleter urlDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request)  {
		const fn = "handlers.url.delete.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("alias not provided"))
			
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("url not found"))
			return
		}

		if err != nil {
			log.Info("failed to delete url", slog.Any("err", err))
			render.JSON(w, r, response.Error("internal error deleting url"))
			return
		}

		log.Info("url deleted", slog.String("alias", alias))
		render.JSON(w, r, response.OK())
	}
}