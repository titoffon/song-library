package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"main/controllers"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Маршрут получение данных библиотеки с фильтрацией по всем полям и пагинацией
	r.Handle("/songs", http.HandlerFunc(controllers.GetSongs)).Methods("GET")

	// Маршрут для получения текста песни с пагинацией по куплетам
	r.HandleFunc("/songs/{id}/verse", controllers.GetSongVerse).Methods("GET")

	// Маршрут для удаления песни
	r.HandleFunc("/songs/{id}", controllers.DeleteSong).Methods("DELETE")

	// Маршрут для обновления песни
	r.HandleFunc("/songs/{id}", controllers.UpdateSong).Methods("PUT")

	//Маршрут для добавления песни
	r.HandleFunc("/songs", controllers.AddSong).Methods("POST")

	return r
}
