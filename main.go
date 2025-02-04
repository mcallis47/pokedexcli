package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(con *config) error
}
type locationPage struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}
type config struct {
	Next     string
	Previous string
}

var commands map[string]cliCommand
var curConfig config

func main() {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		commands = map[string]cliCommand{
			"help": {
				name:        "help",
				description: "Displays a help message",
				callback:    commandHelp,
			},
			"exit": {
				name:        "exit",
				description: "Exit the Pokedex",
				callback:    commandExit,
			},
			"map": {
				name:        "map",
				description: "Displays the names of 20 location areas in the Pokemon world",
				callback:    commandMap,
			},
			"mapb": {
				name:        "mapb",
				description: "Displays the names of 20 location areas in the Pokemon world from the previous page",
				callback:    commandMapB,
			},
		}
		curConfig = config{
			Next:     "",
			Previous: "",
		}
		for {
			fmt.Print("Pokedex > ")
			scanner.Scan()
			text := scanner.Text()

			textCleaned := cleanInput(strings.ToLower(text))
			cmd := textCleaned[0]

			if command, exists := commands[cmd]; exists {
				command.callback(&curConfig)
			} else {
				fmt.Println("Unknown command")
			}
		}
	}

}

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func commandExit(con *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(con *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(con *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if con.Next != "" {
		url = con.Next
	}
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data from the API")
	}
	defer res.Body.Close()
	var locations locationPage
	if err := json.NewDecoder(res.Body).Decode(&locations); err != nil {
		fmt.Println("Error decoding API response")
		return err
	}
	con.Previous = locations.Previous
	con.Next = locations.Next
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(con *config) error {
	if con.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	url := con.Previous
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data from the API")
	}
	defer res.Body.Close()
	var locations locationPage
	if err := json.NewDecoder(res.Body).Decode(&locations); err != nil {
		fmt.Println("Error decoding API response")
		return err
	}
	con.Previous = locations.Previous
	con.Next = locations.Next
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}
