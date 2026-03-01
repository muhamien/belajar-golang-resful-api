package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type CategoryController interface {
	Create(w http.ResponseWriter, r *http.Request, param httprouter.Param)
	Update(w http.ResponseWriter, r *http.Request, param httprouter.Param)
	Delete(w http.ResponseWriter, r *http.Request, param httprouter.Param)
	FindById(w http.ResponseWriter, r *http.Request, param httprouter.Param)
	FindAll(w http.ResponseWriter, r *http.Request, param httprouter.Param)
}
