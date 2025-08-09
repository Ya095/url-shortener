package save

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// omitempty - если пустое поле для параметра - оно просто отсутствует
// иначе пустая строка

// validate - тег, который дает инфо пакету validator/v10

type Request struct {
	URL string `json:"url" validate:"required,url"` // обязательное поле + должен быть валидный url
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	// Сохраняет объект в бд. 
	// Должен копировать сигнатуру storage/sqlite (SaveURL func) + названия должны совпадать
	SaveURL (url string, alias string) (int64, error)
}

// AliasLength - Длина алиаса
const AliasLength = 6


// New - создаёт и возвращает готовый обработчик для сохранения url.
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.save.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		
		// Распарсить запрос и обработать ошибку
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			// Возвращаем ответ с ошибкой клиенту
			render.JSON(w, r, response.Error("failed to decode request"))

			return 
		}

		log.Info("request body decoded", slog.Any("request", req))

		// валидируем структуру req через пакет validator/v10
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(AliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))

			return
		}
		if err != nil {
			log.Info("failed to add url", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		render.JSON(w, r, Response {
			Response: response.OK(),
			Alias: alias,
		})
	}
}