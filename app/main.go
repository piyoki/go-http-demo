package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// just adding space for testing
func main() {
	host := os.Getenv("SERVER_HOST")

	// define the routes handler
	http.HandleFunc("/", HandleGetVideos)
	http.HandleFunc("/update", HandleUpdateVideos)

	http.ListenAndServe(host, nil)
}

// GET
func HandleGetVideos(w http.ResponseWriter, r *http.Request) {

	// fetch json data
	videos := getVideos()

	videoBytes, err := json.Marshal(videos)

	if err != nil {
		panic(err)
	}

	// write back the response (Take in slice of bytes)
	w.Write(videoBytes)
}

// POST
func HandleUpdateVideos(w http.ResponseWriter, r *http.Request) {

	// validate the method
	if r.Method == "POST" {

		// read the body into slice of bytes
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		// create an empty slice of videos
		var videos []video
		// convert the slice of bytes into slice of videos
		err = json.Unmarshal(body, &videos)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Bad request")
		}

		// write back to json file
		saveVideos(videos)

	} else {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not Supported!")
	}
}
