package deleteuser

import (
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

type DeleteUser interface {
	DeleteUser(string) error
}

func New(log *slog.Logger, deleteUser DeleteUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.urlpath.deleteuser.New"

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

		err = deleteUser.DeleteUser(req.PassportNumber)
		if err != nil {
			log.Error("user not exists", sl.Err(err))
			render.JSON(w, r, response.Error("user not deleted"))
			return
		}
		log.Info("user deleted", slog.Any("user", req.PassportNumber))

		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}
}
