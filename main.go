package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "os"
	"path/filepath"
	"text/template"
)

type Movies struct {
	MoviesCount int
	Name        []string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func movieDir() []string {
	var moviesSl []string
	films, err := ioutil.ReadDir("./films")
	check(err)
	for _, film := range films {
		if film.IsDir() {
			continue
		}
		if filepath.Ext(film.Name()) == ".mp4" {
			moviesSl = append(moviesSl, film.Name())
		}
	}
	return moviesSl
}
func viewHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Host, request.RemoteAddr)
	moviesSl := movieDir()
	moviesSt := Movies{
		MoviesCount: len(moviesSl),
		Name:        moviesSl,
	}
	html, err := template.ParseFiles("./index.html")
	check(err)
	err = html.Execute(writer, moviesSt)
	check(err)
}
func viewHandlerVideo(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.FormValue("movie"))
	movie := request.FormValue("movie")
	/*
		moviesSl := movieDir()
		moviesSt := Movies{
			MoviesCount: len(moviesSl),
			Name:        moviesSl,
		}
	*/
	html, err := template.ParseFiles("./template.html")
	check(err)
	err = html.Execute(writer, movie)
	check(err)
}

func main() {

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts/"))))
	http.Handle("/films/", http.StripPrefix("/films/", http.FileServer(http.Dir("./films/"))))
	http.HandleFunc("/", logPanics(viewHandler))
	http.HandleFunc("/video", logPanics(viewHandlerVideo))
	if err := http.ListenAndServe("0.0.0.0:9001", nil); err != nil {
		log.Fatal("failed to start server", err)
	}
}
func logPanics(function func(http.ResponseWriter,
	*http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic: %v", request.RemoteAddr, x)
			}
		}()
		function(writer, request)
	}
}
