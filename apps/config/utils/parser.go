package utils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func ParseID(r *http.Request) (uint, error) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}
