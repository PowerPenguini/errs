package errs

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error) {
	var list *ErrorList
	if first := firstTypedError(err); first != nil {
		list, _ = first.(*ErrorList)
	}
	if list != nil && list.Len() > 0 {
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

		payload := make([]map[string]any, 0, len(list.Errors))
		for _, item := range list.Errors {
			if item == nil {
				continue
			}
			entry := map[string]any{
				"error":   item.Code,
				"message": item.Message,
			}
			if item.Field != "" {
				entry["field"] = item.Field
			}
			payload = append(payload, entry)
		}
		if len(payload) == 0 {
			payload = append(payload, map[string]any{
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
	if first := firstTypedError(err); first != nil {
		e, _ = first.(*Error)
	}
	if e != nil {
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
		payload := map[string]any{
			"error":   e.Code,
			"message": e.Message,
		}
		if e.Field != "" {
			payload["field"] = e.Field
		}
		json.NewEncoder(w).Encode(payload)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "unknown_error",
		"message": "internal server error",
	})
}

// Stop at the first typed boundary. A cause cannot override its parent's status.
func firstTypedError(err error) error {
	switch value := err.(type) {
	case *Error, *ErrorList:
		return err
	case interface{ Unwrap() error }:
		return firstTypedError(value.Unwrap())
	case interface{ Unwrap() []error }:
		for _, cause := range value.Unwrap() {
			if typed := firstTypedError(cause); typed != nil {
				return typed
			}
		}
	}
	return nil
}
