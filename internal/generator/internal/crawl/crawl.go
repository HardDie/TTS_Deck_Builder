package crawl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"tts_deck_build/internal/generator/internal/types"
)

// Parse json file to deck
func parseJson(path string) *types.Deck {
	desc := &types.Deck{}

	// Open file
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Decode
	dec := json.NewDecoder(f)
	if err = dec.Decode(desc); err != nil {
		log.Fatal(err.Error())
	}

	return desc
}

// Check every folder and get cards information
func crawl(path string) (result []*types.Deck) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, file := range files {
		newPath := filepath.Join(path, file.Name())
		if file.IsDir() {
			// Skip 'files' folder
			if file.Name() == "files" {
				continue
			}
			result = append(result, crawl(newPath)...)
			continue
		}
		log.Println("Parse file:", newPath)

		// Parse only json files
		if filepath.Ext(newPath) != ".json" {
			continue
		}
		deck := parseJson(newPath)

		result = append(result, deck)
		// Set for each card
		for _, card := range deck.Cards {
			card.FillWithInfo(deck.Version, deck.Collection, deck.Type)
		}
	}
	return
}

// Separate decks by type
// Top level map[string] - split by types (ex.: Loot, Monster)
// Next level []*Deck - split by collection (ex.: Base, DLC)
func Crawl(path string) map[string][]*types.Deck {
	result := make(map[string][]*types.Deck)
	// Get all decks
	decks := crawl(path)
	// Split decks by type
	for _, deck := range decks {
		result[deck.Type] = append(result[deck.Type], deck)
	}
	return result
}