package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var (
	students []Student
	idSeq    = 1
	mu       sync.Mutex
)

func getStudents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(students)
	log.Println("Get all students: ")
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	log.Println("getStudent: ", id)
	for _, s := range students {
		if s.ID == id {
			json.NewEncoder(w).Encode(s)
			return
		}
	}
	http.NotFound(w, r)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	var s Student
	json.NewDecoder(r.Body).Decode(&s)
	log.Println("createStudent: ")
	mu.Lock()
	s.ID = idSeq
	idSeq++
	students = append(students, s)
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	log.Println("updateStudent: ", id)
	var updated Student
	json.NewDecoder(r.Body).Decode(&updated)

	mu.Lock()
	defer mu.Unlock()
	for i, s := range students {
		if s.ID == id {
			students[i].Name = updated.Name
			students[i].Age = updated.Age
			json.NewEncoder(w).Encode(students[i])
			return
		}
	}
	http.NotFound(w, r)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)
	log.Println("deleteStudent: ", id)
	mu.Lock()
	defer mu.Unlock()
	for i, s := range students {
		if s.ID == id {
			students = append(students[:i], students[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudent).Methods("GET")
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")

	port := "10000"
	log.Println("Server is running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
