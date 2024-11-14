package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleError(err error, w http.ResponseWriter) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(
			fmt.Sprintf("ERROR: %q", err),
		))
	}
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the to do app"))
	w.Write([]byte("Available endpoints are create/ read/ update/ delete/"))
}

type createHandler struct {
	dal DataAccessLayer
}

func (h *createHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data ToDoItem
		json.NewDecoder(r.Body).Decode(&data)
		err := h.dal.Create(data)
		handleError(err, w)
	}
}

type readHandler struct {
	dal DataAccessLayer
}

func (h *readHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	items := h.dal.Read()
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		errorMsg := fmt.Sprintf("ERROR! %q", err)
		fmt.Println(errorMsg)
		w.Write([]byte(errorMsg))
	} else {
		w.Write([]byte(data))
	}
}

type updateHandler struct {
	dal DataAccessLayer
}

func (h *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data ToDoItem
		json.NewDecoder(r.Body).Decode(&data)
		err := h.dal.Update(data)
		handleError(err, w)
	}
}

type deleteHandler struct {
	dal DataAccessLayer
}

func (h *deleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data ToDoItem
		json.NewDecoder(r.Body).Decode(&data)
		err := h.dal.Delete(data)
		handleError(err, w)
	}
}

func StartAPI(dal DataAccessLayer) {
	mux := http.NewServeMux()
	mux.Handle("/", &homeHandler{})
	mux.Handle("/create", &createHandler{dal})
	mux.Handle("/read", &readHandler{dal})
	mux.Handle("/update", &updateHandler{dal})
	mux.Handle("/delete", &deleteHandler{dal})
	http.ListenAndServe(":8080", mux)
}
