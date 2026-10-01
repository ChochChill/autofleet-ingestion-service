package main

import (
	"carsAPI/handler"
	"carsAPI/repository"
	"log"
	"net/http"
	"carsAPI/db"
)

func main() {
	database, err := db.OpenDatabase()
	if err != nil{
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.CreateTables(database); err != nil {
		log.Fatal(err)
	}
	repo := repository.NewCarRepository(database)
	carHandler := handler.NewCarHandler(repo)
	http.HandleFunc("GET /cars", carHandler.GetCars)
	http.HandleFunc("GET /cars/{id}", carHandler.GetCar)
	http.HandleFunc("POST /cars", carHandler.AddCar)
	http.HandleFunc("PUT /cars/{id}", carHandler.UpdateCar)
	http.HandleFunc("DELETE /cars/{id}", carHandler.DeleteCar)

	log.Println("carsinfo server starting on :3015")
	log.Fatal(http.ListenAndServe(":3015", nil))
}
