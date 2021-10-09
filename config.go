package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	MapArray       []int   `json:"mapArray"`
	MapX           int     `json:"mapX"`
	MapY           int     `json:"mapY"`
	SpawnLocationX float64 `json:"spawnLocationX"`
	SpawnLocationY float64 `json:"spawnLocationY"`
}

func config() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	
	mapX = configuration.MapX
	mapY = configuration.MapY
	map_array = configuration.MapArray
	player_pos_x = configuration.SpawnLocationX
	player_pos_y = configuration.SpawnLocationY
}
