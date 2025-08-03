package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/ShamirMaharjan/Student-management-system--GO/internal/body"
	"github.com/ShamirMaharjan/Student-management-system--GO/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func CreateStudent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating Student")

		var student body.Student

		//decode returns an error if the request body is empty
		err := json.NewDecoder(r.Body).Decode(&student)
		//EOF is when nothing is sent in the request body
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request body validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}
