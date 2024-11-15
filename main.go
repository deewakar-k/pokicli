package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Pokemon struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Types  []Type `json:"types"`
	Stats  []Stat `json:"stats"`
}

type Type struct {
	Slot int       `json:"slot"`
	Type TypeValue `json:"type"`
}

type TypeValue struct {
	Name string `json:"name"`
}

type Stat struct {
	BaseStat    int      `json:"base_stat"`
	StatDetails StatInfo `json:"stat"`
}

type StatInfo struct {
	Name string `json:"name"`
}

func main() {

	pokiName := os.Args[1]
	api := "https://pokeapi.co/api/v2/pokemon/" + pokiName
	resp, err := http.Get(api)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var poki Pokemon

	err = json.Unmarshal(body, &poki)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("name: %s\n", poki.Name)
	fmt.Println("types:")
	for _, t := range poki.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	fmt.Println("stats:")
	for _, stat := range poki.Stats {
		fmt.Printf("  - %s: %d\n", stat.StatDetails.Name, stat.BaseStat)
	}
}
