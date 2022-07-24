package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sanber-rest-api-mysql/config"
	"sanber-rest-api-mysql/models"
	"sanber-rest-api-mysql/movie"
	"sanber-rest-api-mysql/utils"

	"github.com/julienschmidt/httprouter"
)

func main() {
	db, e := config.MySQL()

	if e != nil {
		log.Fatal(e)
	}

	eb := db.Ping()
	if eb != nil {
		panic(eb.Error())
	}

	fmt.Println("Success")

	router := httprouter.New()

	// Endpoint
	router.GET("/movie", GetMovie)
	router.POST("/movie/create", PostMovie)
	router.PUT("/movie/:id/update", UpdateMovie)
	router.DELETE("/movie/:id/delete", DeleteMovie)

	fmt.Println("Server Running at Port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

// Get Data Movie
func GetMovie(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	movies, err := movie.GetAll(ctx)

	if err != nil {
		fmt.Println(err)
	}

	utils.ResponseJSON(w, movies, http.StatusOK)
}

// Create Data Movie
func PostMovie(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(
			w,
			"Use Content-Type - application/json",
			http.StatusBadRequest,
		)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mov models.Movie
	if err := json.NewDecoder(r.Body).Decode(&mov); err != nil {
		utils.ResponseJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := movie.Insert(ctx, mov); err != nil {
		utils.ResponseJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := map[string]string{
		"status": "Success",
	}
	utils.ResponseJSON(w, res, http.StatusCreated)

}

// Update Data Movie
func UpdateMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(
			w,
			"Use Content-Type - application/json",
			http.StatusBadRequest,
		)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var mov models.Movie

	if err := json.NewDecoder(r.Body).Decode(&mov); err != nil {
		utils.ResponseJSON(w, err, http.StatusBadRequest)
		return
	}

	var idMovie = ps.ByName("id")

	if err := movie.Update(ctx, mov, idMovie); err != nil {
		utils.ResponseJSON(w, err, http.StatusInternalServerError)
		return
	}

	res := map[string]string{
		"status": "Success",
	}
	utils.ResponseJSON(w, res, http.StatusCreated)

}

// Dalate Data Movie
func DeleteMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var idMovie = ps.ByName("id")

	if err := movie.Delete(ctx, idMovie); err != nil {
		errorMessage := map[string]string{
			"error": fmt.Sprintf("%v", err),
		}
		utils.ResponseJSON(w, errorMessage, http.StatusInternalServerError)
		return
	}

	res := map[string]string{
		"status": "Success",
	}
	utils.ResponseJSON(w, res, http.StatusOK)

}
