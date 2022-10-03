package config

import "path/filepath"

const (
	MaxFilenameLength = 200

	MinWidth  = 2
	MinHeight = 2
	MaxWidth  = 10
	MaxHeight = 7
	MaxCount  = MaxWidth*MaxHeight - 1
)

type Config struct {
	Debug bool `json:"debug"`

	Data   string `json:"data"`
	Game   string `json:"game"`
	Cache  string `json:"cache"`
	Result string `json:"result"`

	Setting       string `json:"settings"`
	InfoFilename  string `json:"infoFilename"`
	ImageFilename string `json:"imageFilename"`

	CardImagePath       string `json:"cardImagePath"`
	DeckImagePath       string `json:"deckImagePath"`
	CollectionImagePath string `json:"collectionImagePath"`
	GameImagePath       string `json:"gameImagePath"`
}

func Get() *Config {
	return &Config{
		Debug: false,

		Data:   "data",
		Game:   "games",
		Cache:  "cache",
		Result: "result",

		Setting:       "settings.json",
		InfoFilename:  ".info.json",
		ImageFilename: ".image.bin",

		CardImagePath:       "/api/games/%s/collections/%s/decks/%s/cards/%d/image",
		DeckImagePath:       "/api/games/%s/collections/%s/decks/%s/image",
		CollectionImagePath: "/api/games/%s/collections/%s/image",
		GameImagePath:       "/api/games/%s/image",
	}
}

func (c *Config) Settings() string {
	return filepath.Join(c.Data, c.Setting)
}
func (c *Config) Games() string {
	return filepath.Join(c.Data, c.Game)
}
func (c *Config) Caches() string {
	return filepath.Join(c.Data, c.Cache)
}
func (c *Config) Results() string {
	return filepath.Join(c.Data, c.Result)
}

// SetDataPath For tests only!!!
func (c *Config) SetDataPath(dataPath string) {
	c.Data = dataPath
}
