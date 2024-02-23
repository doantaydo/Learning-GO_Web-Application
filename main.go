package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Animes struct {
	AnimeName   string `json:"anime_name"`
	ReleaseYear int    `json:"release_year"`
	NumOfEp     int    `json:"num_of_ep"`
}

func main() {
	user1 := Animes{
		AnimeName:   "A",
		ReleaseYear: 2,
		NumOfEp:     1,
	}

	user2 := Animes{
		AnimeName:   "A",
		ReleaseYear: 2,
		NumOfEp:     1,
	}

	var newAnimeList []Animes
	newAnimeList = append(newAnimeList, user1)
	newAnimeList = append(newAnimeList, user2)

	newJson, err := json.MarshalIndent(newAnimeList, "", "    ")

	if err != nil {
		fmt.Println("Mershal err: ", err)
	}

	log.Println(string(newJson))
}
