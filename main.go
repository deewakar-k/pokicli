package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Pokemon struct {
	Name    string `json:"name"`
	Weight  int    `json:"weight"`
	Types   []Type `json:"types"`
	Stats   []Stat `json:"stats"`
	Sprites Sprite `json:"sprites"`
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

type Sprite struct {
	FrontDefault string `json:"front_default"`
}

func downloadSprite(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("please provide a pokemon name!")
	}

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

	sprite, err := downloadSprite(poki.Sprites.FrontDefault)
	if err != nil {
		log.Fatalln("failed to download image: ", err)
	}

	cmd := exec.Command("kitty", "+kitten", "icat")
	cmd.Stdin = bytes.NewReader(sprite)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("failed to display image: %v", err)
	}

	fmt.Printf("\nname: %s\n", poki.Name)
	fmt.Println("types:")
	for _, t := range poki.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	fmt.Println("stats:")
	for _, stat := range poki.Stats {
		fmt.Printf("  - %s: %d\n", stat.StatDetails.Name, stat.BaseStat)
	}
}
