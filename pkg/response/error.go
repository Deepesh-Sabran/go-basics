package response

import (
	"encoding/json"
	"net/http"

	"github.com/Deepesh-Sabran/go-basics/internal/errors"
)

func HandleError(w http.ResponseWriter, err error) {
	appErr, ok:= err.(*errors.AppError)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "internal server error",
			"code": "INTERNAL_SERVER_ERROR",
		})

		return
	}

	w.WriteHeader(appErr.StatusCode)
	json.NewEncoder(w).Encode(appErr)
}