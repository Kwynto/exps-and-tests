package user

import (
	"net/http"
	"rest-api/internal/handlers"
	"rest-api/pkg/simlog"
)

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc(usersURL+"/list", h.GetList)
	router.HandleFunc(usersURL+"/create", h.CreateUser)
	router.HandleFunc(usersURL+"/uuid", h.GetUserByUUID)
	router.HandleFunc(usersURL+"/update", h.UpdateUser)
	router.HandleFunc(usersURL+"/partupdate", h.PartiallyUpdateUser)
	router.HandleFunc(usersURL+"/delete", h.DeleteUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		simlog.Debug("Faild to get list of users")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Faild to get list of users. Method is not allowed."))
	}
	simlog.Debug("get list of users")
	w.WriteHeader(http.StatusOK)             // 200
	w.Write([]byte("this is list of users")) // Заглушка
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) // 201
	w.Write([]byte("Create user"))    // Заглушка
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)           // 200
	w.Write([]byte("this is get of user")) // Заглушка
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)       // 204
	w.Write([]byte("this is update of user")) // Заглушка
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)                 // 204
	w.Write([]byte("this is partially update of user")) // Заглушка
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)       // 204
	w.Write([]byte("this is delete of user")) // Заглушка
}
