package main

import (
	"encoding/json"
	"io/ioutil"
)

type video struct {
	Id          string
	Title       string
	Description string
	Imageurl    string
	Url         string
}

func getVideos() (videos []video) {

	fileBytes, err := ioutil.ReadFile("./videos.json")

	if err != nil {
		panic(err)
	}

	// convert the slice of bytes into a slices
	err = json.Unmarshal(fileBytes, &videos)

	if err != nil {
		panic(err)
	}

	return videos
}

func saveVideos(videos []video) {

	// convert back into slice of bytes
	videoBytes, err := json.Marshal(videos)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./videos-updated.json", videoBytes, 0644)
	if err != nil {
		panic(err)
	}

}
