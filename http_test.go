package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteErrorPreservesErrorLists(t *testing.T) {
	list := NewErrorList(NewFieldError("name_required", "name", "Name is required.", ValidationType, nil), NewFieldError("age_invalid", "age", "Age is invalid.", ValidationType, nil))
	for name, err := range map[string]error{"direct": list, "wrapped": fmt.Errorf("context: %w", list), "joined": errors.Join(errors.New("context"), list)} {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteError(w, err)
			if w.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
			}
			var body struct {
				Errors []struct{ Error, Field, Message string }
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if len(body.Errors) != 2 || body.Errors[0].Error != "name_required" || body.Errors[0].Field != "name" || body.Errors[0].Message != "Name is required." || body.Errors[1].Field != "age" {
				t.Fatalf("list contract lost: %+v", body)
			}
		})
	}
}

func TestWriteErrorUsesOuterTypedBoundary(t *testing.T) {
	list := NewErrorList(NewFieldError("name_required", "name", "Name is required.", ValidationType, nil))
	for typ, status := range map[string]int{BadRequestType: 400, ValidationType: 422, InternalType: 500, NotFoundType: 404, UnauthorizedType: 401, ForbiddenType: 403} {
		t.Run(typ, func(t *testing.T) {
			outer := NewFieldError("outer_error", "resource", "Outer message.", typ, list)
			for _, err := range []error{outer, fmt.Errorf("context: %w", outer), errors.Join(errors.New("context"), outer)} {
				w := httptest.NewRecorder()
				WriteError(w, err)
				var body map[string]any
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if w.Code != status || body["error"] != outer.Code || body["message"] != outer.Message || body["field"] != outer.Field || body["errors"] != nil {
					t.Fatalf("cause overrode boundary: status=%d body=%s", w.Code, w.Body.String())
				}
				if w.Header().Get("Content-Type") != "application/json" {
					t.Fatal("missing JSON content type")
				}
			}
		})
	}
}

func TestWriteErrorDoesNotExposeUntypedCause(t *testing.T) {
	var typedNil *Error
	var listNil *ErrorList
	for name, err := range map[string]error{"raw": errors.New("SQL password diagnostic"), "nil": nil, "typed_nil": typedNil, "list_nil": listNil, "empty_list": NewErrorList()} {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteError(w, err)
			var body map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != 500 || body["error"] != "unknown_error" || body["message"] != "internal server error" || strings.Contains(w.Body.String(), "SQL") {
				t.Fatalf("unsafe fallback: status=%d body=%s", w.Code, w.Body.String())
			}
			if w.Header().Get("Content-Type") != "application/json" {
				t.Fatal("missing JSON content type")
			}
		})
	}
}
