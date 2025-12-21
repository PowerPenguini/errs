package errs

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error) {
	var list *ErrorList
	if errors.As(err, &list) && list != nil && list.Len() > 0 {
		status := http.StatusBadRequest
		for _, item := range list.Errors {
			if item == nil {
				continue
			}
			if item.Type == ValidationType {
				status = http.StatusUnprocessableEntity
				break
			}
			if item.Type == InternalType {
				status = http.StatusInternalServerError
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		payload := make([]map[string]string, 0, len(list.Errors))
		for _, item := range list.Errors {
			if item == nil {
				continue
			}
			payload = append(payload, map[string]string{
				"error":   item.Code,
				"message": item.Message,
			})
		}
		if len(payload) == 0 {
			payload = append(payload, map[string]string{
				"error":   "unknown_error",
				"message": "unknown error",
			})
		}

		json.NewEncoder(w).Encode(map[string]any{
			"errors": payload,
		})
		return
	}

	var e *Error
	if errors.As(err, &e) {
		status := http.StatusInternalServerError
		switch e.Type {
		case ValidationType:
			status = http.StatusUnprocessableEntity
		case InternalType:
			status = http.StatusInternalServerError
		case NotFoundType:
			status = http.StatusNotFound
		case UnauthorizedType:
			status = http.StatusUnauthorized
		case ForbiddenType:
			status = http.StatusForbidden
		default:
			status = http.StatusBadRequest
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   e.Code,
			"message": e.Message,
		})
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "unknown_error",
		"message": err.Error(),
	})
}
