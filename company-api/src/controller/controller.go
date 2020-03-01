package controller

import (
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var (
	validate = validator.New()
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
