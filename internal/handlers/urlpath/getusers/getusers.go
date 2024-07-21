package getusers

import (
	"errors"
	"log/slog"
	"net/http"

	"timetracker_effective_mobile/internal/lib/api/response"
	"timetracker_effective_mobile/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	PassportNumber string `json:"passportNumber" validate:"required,url"`
}

type Response struct {
	response.Response
}

type UserSaver interface {
	AddUser(PassportNumber string) (id int, err error)
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.urlpath.adduser.New"

		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("can't decode request", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if len(req.PassportNumber) != 11 {
			log.Error("invalid passport number", sl.Err(err))
			render.JSON(w, r, response.Error("invalid passport number"))
			return
		}

		id, err := userSaver.AddUser(req.PassportNumber)
		if errors.Is(err, errors.New("user exists")) {
			log.Error("user exists", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add user"))
			return
		}

		if err != nil {
			log.Error("can't add user", sl.Err(err))
			render.JSON(w, r, response.Error("failed to add user"))
			return
		}
		log.Info("user added", slog.Int("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
