package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if params["id"] == item.ID {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	new := json.NewDecoder(r.Body).Decode(&movie)
	fmt.Println("Entered Movie", new)
	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, item := range movies {
		if params["id"] == item.ID {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			new := json.NewDecoder(r.Body).Decode(&movie)
			fmt.Println("Entered Movie", new)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "978-0345339465", Title: "Titanic", Director: &Director{Firstname: "James", Lastname: "Cameron"}})
	movies = append(movies, Movie{ID: "2", Isbn: "978-0345387465", Title: "The GodFather", Director: &Director{Firstname: "Francis Ford", Lastname: "Coppola"}})

	routes := []map[string]interface{}{
		{"path": "/movies", "method": "GET", "handlerFunctionName": getMovies, "description": "Get all movies"},
		{"path": "/movie/{id}", "method": "DELETE", "handlerFunctionName": deleteMovie, "description": "Delete a movie by id"},
		{"path": "/movie", "method": "POST", "handlerFunctionName": createMovie, "description": "Create new movie"},
		{"path": "/movie/{id}", "method": "PUT", "handlerFunctionName": updateMovie, "description": "Update a movie by id"},
		{"path": "/movie/{id}", "method": "GET", "handlerFunctionName": getMovie, "description": "Get single movie by id"},
	}

	for _, route := range routes {
		r.HandleFunc(route["path"].(string), route["handlerFunctionName"].(func(http.ResponseWriter, *http.Request))).Methods(route["method"].(string))
	}

	fmt.Printf("Starting Server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
