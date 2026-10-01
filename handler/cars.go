package handler

import (
	"carsAPI/model"
	"carsAPI/repository"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CarHandler struct {
	repo *repository.CarRepository
}

func NewCarHandler(repo *repository.CarRepository) *CarHandler {
	return &CarHandler{
		repo: repo,
	}
}
func (h *CarHandler) GetCars(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	cars, err := h.repo.GetAll()
	if err != nil {
		logAndReturnError(err, 0, "failed to get cars", w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func (h *CarHandler) GetCar(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		logAndReturnError(err, 0, "invalid ID", w, http.StatusBadRequest)
		return
	}
	car, found, err := h.repo.GetByID(id)
	if err != nil {
		logAndReturnError(err, id, "failed to get car", w, http.StatusInternalServerError)
		return
	}
	if found {
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(car)
		return
	}
	logAndReturnError(nil, id, "car not found", w, http.StatusNotFound)
}

func (h *CarHandler) AddCar(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	var car *model.Car
	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		logAndReturnError(err, 0, "invalid body sent", w, http.StatusBadRequest)
		return
	}
	car, err = h.repo.Create(car)
	if err != nil {
		logAndReturnError(err, 0, "add car failed", w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
	logRequestCompletion(r.Method, car.ID)
}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		logAndReturnError(err, 0, "invalid ID", w, http.StatusBadRequest)
		return
	}
	isDeleted, err := h.repo.Delete(id)
	if err != nil {
		logAndReturnError(err, 0, "delete failed", w, http.StatusInternalServerError)
		return
	}
	if isDeleted {
		w.WriteHeader(http.StatusNoContent)
		logRequestCompletion(r.Method, id)
		return
	}
	logAndReturnError(nil, id, "car not found", w, http.StatusNotFound)
}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	var newCar *model.Car
	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		logAndReturnError(err, 0, "invalid body sent", w, http.StatusBadRequest)
		return
	}

	idParam := r.PathValue("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		logAndReturnError(err, 0, "invalid ID", w, http.StatusBadRequest)
		return
	}
	newCar, isUpdated, err := h.repo.Update(id, newCar)
	if err != nil {
		logAndReturnError(err, id, "update failed", w, http.StatusInternalServerError)
		return
	}
	if isUpdated {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newCar)
		logRequestCompletion(r.Method, id)
		return
	}
	logAndReturnError(nil, id, "car not found", w, http.StatusNotFound)
}

func logRequest(r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
}
func logRequestCompletion(methodType string, id int) {
	log.Printf("%s done for car with ID:%d", methodType, id)
}
func logAndReturnError(err error, id int, message string, w http.ResponseWriter, resStatus int) {
	log.Printf("%s, ID:%d, error: %v", message, id, err)
	http.Error(w, message, resStatus)
}
