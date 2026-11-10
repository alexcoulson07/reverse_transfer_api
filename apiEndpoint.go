package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/alexcoulson07/reversetransfer/api/apiFunctions"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func transferCourses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	subject := q.Get("subject")
	number := q.Get("number")
	if subject == "" || number == "" {
		http.Error(w, "Bad parameters", http.StatusBadRequest)
		return
	}
	data, err := apiFunctions.GetTransferCourses(subject, number)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Courses apiFunctions.DetailedCourses `json:"transfer_courses"`
	}{
		Courses: data,
	}
	json.NewEncoder(w).Encode(response)
}

func subjects(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sgbi := q.Get("sgbi")
	if sgbi == "" {
		http.Error(w, "Bad parameters", http.StatusBadRequest)
		return
	}
	data, err := apiFunctions.GetSubjects(sgbi)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Subjects []string `json:"subjects"`
	}{
		Subjects: data,
	}
	json.NewEncoder(w).Encode(response)
}

func courses(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sgbi := q.Get("sgbi")
	subject := q.Get("subject")
	if sgbi == "" || subject == "" {
		http.Error(w, "Bad parameters", http.StatusBadRequest)
		return
	}
	data, err := apiFunctions.GetCourses(sgbi, subject)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Courses apiFunctions.Courses `json:"school_courses"`
	}{
		Courses: data,
	}
	json.NewEncoder(w).Encode(response)
}

func search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	sgbi := q.Get("sgbi")
	search := q.Get("search")
	if sgbi == "" || search == "" {
		http.Error(w, "Bad parameters", http.StatusBadRequest)
		return
	}
	data, err := apiFunctions.SearchCourses(sgbi, search)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	response := struct {
		Result apiFunctions.SearchResults `json:"result"`
	}{
		Result: data,
	}
	json.NewEncoder(w).Encode(response)
}

func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/api/transferCourses", transferCourses)
	router.HandleFunc("/api/schoolSubjects", subjects)
	router.HandleFunc("/api/subjectCourses", courses)
	router.HandleFunc("/api/search", search)
	log.Fatal(http.ListenAndServe(":4322", handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:4321"}),
		handlers.AllowedMethods([]string{"GET"}),
	)(router)))

}

func main() {
	fmt.Println(apiFunctions.InitDb())
	handleRequest()

}
