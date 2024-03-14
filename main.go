package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"code.sajari.com/docconv"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Resp struct {
	UniqueName string `json:"UniqueName"`
	Body       string `json:"Body"`
}

func ConvFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["id"]

	res, err := docconv.ConvertPath(name)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(name)
	io.WriteString(w, res.Body)
}

func ConvNewFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file attached", http.StatusBadRequest)
		return
	}
	defer file.Close()

	nameString := fileHeader.Filename
	ext := strings.Split(nameString, ".")

	name := uuid.NewString()
	dst, err := os.Create("/app/files/" + name + "." + ext[len(ext)-1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	fmt.Println(dst.Name())

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := docconv.ConvertPath(dst.Name())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(Resp{name, res.Body})

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(result))
}

func initRouter() {

	r := mux.NewRouter()
	r.HandleFunc("/{id}", ConvFile).Methods("GET")
	r.HandleFunc("/newfile", ConvNewFile).Methods("POST")

	fmt.Println("started")
	err := http.ListenAndServe("0.0.0.0:3334", r)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	initRouter()
}
