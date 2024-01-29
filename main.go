package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
	//director must be a pointer to Director struct, because it is a nested struct
	// and we need to access the fields of the nested struct
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	//Direct not exacltly represent a db entity, it is just a nested struct
}

var movies []Movie //slice of Movie struct

func PopulateMovies() {
	movies = append(movies,
		Movie{Id: "1", Isbn: "123", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies,
		Movie{Id: "2", Isbn: "456", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	movies = append(movies,
		Movie{Id: "3", Isbn: "789", Title: "Movie Three", Director: &Director{Firstname: "Jane", Lastname: "Doe"}})
}

func Pong(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode("Pong")
}

func GetMovies(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(movies)
}

func GetMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	id := vars["id"]
	for _, movie := range movies {
		if movie.Id == id {
			json.NewEncoder(response).Encode(movie)
			return
		}
	}
	json.NewEncoder(response).Encode("Movie not found")
}

func DeleteMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	id := vars["id"]
	for idx, movie := range movies {
		if movie.Id == id {
			movies = append(movies[:idx], movies[idx+1:]...)
			json.NewEncoder(response).Encode(
				"Movie with id " + id + " has been deleted")
			return
		}
	}
	json.NewEncoder(response).Encode("Movie not found")
}

func CreateMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var movie Movie
	r := json.NewDecoder(request.Body)
	err := r.Decode(&movie)
	if err != nil {
		fmt.Println(err)
		return
	}
	movie.Id = fmt.Sprintf("%d", len(movies)+1)
	movies = append(movies, movie)
	json.NewEncoder(response).Encode("Movie has been created")
}

func UpdateMovie(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(request)
	id := vars["id"]

	for idx, movie := range movies {
		if movie.Id == id {
			var movie Movie
			r := json.NewDecoder(request.Body)
			err := r.Decode(&movie)
			if err != nil {
				fmt.Println(err)
				return
			}
			movie.Id = id
			movies[idx] = movie
			json.NewEncoder(response).Encode("Movie has been updated")
			return
		}
	}

	json.NewEncoder(response).Encode("Movie not found")
}

func main() {
	fmt.Println("Starting the application...")
	PopulateMovies()
	fmt.Printf("Movies has been populated with %d movies\n", len(movies))
	r := mux.NewRouter()
	r.HandleFunc("/ping", Pong).Methods("GET")
	r.HandleFunc("/movies", GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", GetMovie).Methods("GET")
	r.HandleFunc("/movies", CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", DeleteMovie).Methods("DELETE")
	fmt.Printf("Starting app on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
