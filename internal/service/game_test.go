package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	"tts_deck_build/internal/config"
	"tts_deck_build/internal/dto"
	"tts_deck_build/internal/entity"
	er "tts_deck_build/internal/errors"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/utils"
)

type gameTest struct {
	cfg         *config.Config
	gameService IGameService
}

func newGameTest(dataPath string) *gameTest {
	cfg := config.GetConfig()
	cfg.SetDataPath(dataPath)

	return &gameTest{
		cfg:         cfg,
		gameService: NewGameService(repository.NewGameRepository(cfg)),
	}
}

func (tt *gameTest) testCreate(t *testing.T) {
	gameName := "one"
	desc := "best game ever"

	// Create game
	game, err := tt.gameService.Create(&dto.CreateGameDTO{
		Name:        gameName,
		Description: desc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name.String() != gameName {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName)
	}
	if game.Description.String() != desc {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc)
	}

	// Try to create duplicate
	_, err = tt.gameService.Create(&dto.CreateGameDTO{
		Name: gameName,
	})
	if err == nil {
		t.Fatal("Error, you can't create duplicate game")
	}
	if !errors.Is(err, er.GameExist) {
		t.Fatal(err)
	}

	// Delete game
	err = tt.gameService.Delete(game.ID)
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *gameTest) testDelete(t *testing.T) {
	gameName := "one"
	gameID := utils.NameToID(gameName)

	// Try to remove non-existing game
	err := tt.gameService.Delete(gameID)
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.gameService.Create(&dto.CreateGameDTO{
		Name: gameName,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Delete game
	err = tt.gameService.Delete(gameID)
	if err != nil {
		t.Fatal(err)
	}

	// Try to delete game twice
	err = tt.gameService.Delete(gameID)
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}
}
func (tt *gameTest) testUpdate(t *testing.T) {
	gameName := []string{"one", "two"}
	desc := []string{"first description", "second description"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to update non-existing game
	_, err := tt.gameService.Update(gameID[0], &dto.UpdateGameDTO{})
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	game, err := tt.gameService.Create(&dto.CreateGameDTO{
		Name:        gameName[0],
		Description: desc[0],
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name.String() != gameName[0] {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName[0])
	}
	if game.Description.String() != desc[0] {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc[0])
	}

	// Update game
	game, err = tt.gameService.Update(gameID[0], &dto.UpdateGameDTO{
		Name:        gameName[1],
		Description: desc[1],
	})
	if err != nil {
		t.Fatal(err)
	}
	if game.Name.String() != gameName[1] {
		t.Fatal("Bad name [got]", game.Name, "[want]", gameName[1])
	}
	if game.Description.String() != desc[1] {
		t.Fatal("Bad description [got]", game.Description, "[want]", desc[1])
	}

	// Delete game
	err = tt.gameService.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Try to update non-existing game
	_, err = tt.gameService.Update(gameID[1], &dto.UpdateGameDTO{})
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}
}
func (tt *gameTest) testList(t *testing.T) {
	gameName := []string{"B game", "A game"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Empty list
	items, err := tt.gameService.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}

	// Create first game
	_, err = tt.gameService.Create(&dto.CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// One game
	items, err = tt.gameService.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatal("List should be with 1 element")
	}

	// Create second game
	_, err = tt.gameService.Create(&dto.CreateGameDTO{
		Name: gameName[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Sort by name
	items, err = tt.gameService.List("name")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[0])
	}

	// Sort by name_desc
	items, err = tt.gameService.List("name_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[1])
	}

	// Sort by created date
	items, err = tt.gameService.List("created")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[0])
	}
	if items[1].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[1])
	}

	// Sort by created_desc
	items, err = tt.gameService.List("created_desc")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatal("List should with 2 value")
	}
	if items[0].Name.String() != gameName[1] {
		t.Fatal("Bad name order: [got]", items[0].Name, "[want]", gameName[1])
	}
	if items[1].Name.String() != gameName[0] {
		t.Fatal("Bad name order: [got]", items[1].Name, "[want]", gameName[0])
	}

	// Delete first game
	err = tt.gameService.Delete(gameID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Delete second game
	err = tt.gameService.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Empty list
	items, err = tt.gameService.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Fatal("List should be empty")
	}
}
func (tt *gameTest) testItem(t *testing.T) {
	gameName := []string{"one", "two"}
	gameID := []string{utils.NameToID(gameName[0]), utils.NameToID(gameName[1])}

	// Try to get non-existing game
	_, err := tt.gameService.Item(gameID[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.gameService.Create(&dto.CreateGameDTO{
		Name: gameName[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = tt.gameService.Item(gameID[0])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = tt.gameService.Item(gameID[1])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Rename game
	_, err = tt.gameService.Update(gameID[0], &dto.UpdateGameDTO{Name: gameName[1]})
	if err != nil {
		t.Fatal(err)
	}

	// Get valid game
	_, err = tt.gameService.Item(gameID[1])
	if err != nil {
		t.Fatal(err)
	}

	// Get invalid game
	_, err = tt.gameService.Item(gameID[0])
	if err == nil {
		t.Fatal("Error, game not exist")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = tt.gameService.Delete(gameID[1])
	if err != nil {
		t.Fatal(err)
	}
}
func (tt *gameTest) testImage(t *testing.T) {
	gameName := "one"
	gameID := utils.NameToID(gameName)
	pngImage := "https://github.com/fluidicon.png"
	jpegImage := "https://avatars.githubusercontent.com/apple"

	// Check no game
	_, _, err := tt.gameService.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game not exists")
	}
	if !errors.Is(err, er.GameNotExists) {
		t.Fatal(err)
	}

	// Create game
	_, err = tt.gameService.Create(&dto.CreateGameDTO{
		Name:  gameName,
		Image: pngImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err := tt.gameService.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "png" {
		t.Fatal("Image type error! [got]", imgType, "[want] png")
	}

	// Update game
	_, err = tt.gameService.Update(gameID, &dto.UpdateGameDTO{
		Image: jpegImage,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check image type
	_, imgType, err = tt.gameService.GetImage(gameID)
	if err != nil {
		t.Fatal(err)
	}
	if imgType != "jpeg" {
		t.Fatal("Image type error! [got]", imgType, "[want] jpeg")
	}

	// Update game
	_, err = tt.gameService.Update(gameID, &dto.UpdateGameDTO{
		Image: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check no image
	_, _, err = tt.gameService.GetImage(gameID)
	if err == nil {
		t.Fatal("Error, game don't have image")
	}
	if !errors.Is(err, er.GameImageNotExists) {
		t.Fatal(err)
	}

	// Delete game
	err = tt.gameService.Delete(gameID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGame(t *testing.T) {
	t.Parallel()

	// Set path for the game test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		t.Fatal("TEST_DATA_PATH must be set")
	}
	tt := newGameTest(filepath.Join(dataPath, "game_test"))

	t.Run("create", tt.testCreate)
	t.Run("delete", tt.testDelete)
	t.Run("update", tt.testUpdate)
	t.Run("list", tt.testList)
	t.Run("item", tt.testItem)
	t.Run("image", tt.testImage)
}

func (tt *gameTest) fuzzCleanup() {
	_ = os.RemoveAll(tt.cfg.Data)
}
func (tt *gameTest) fuzzList(t *testing.T, waitItems int) error {
	items, err := tt.gameService.List("")
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if len(items) != waitItems {
		{
			data, _ := json.MarshalIndent(items, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("items: [wait] %d, [got] %d", waitItems, len(items))
	}
	return nil
}
func (tt *gameTest) fuzzItem(t *testing.T, gameID, name, desc string) error {
	game, err := tt.gameService.Item(gameID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	if game.Name.String() != name {
		{
			data, _ := json.MarshalIndent(game, "", "	")
			t.Log(string(data))
		}
		return fmt.Errorf("name: [wait] %s [got] %s", name, game.Name)
	}
	if game.Description.String() != desc {
		{
			data, _ := json.MarshalIndent(game, "", "	")
			t.Log("item:", string(data))
		}
		return fmt.Errorf("description: [wait] %q [got] %q", desc, game.Description)
	}
	return nil
}
func (tt *gameTest) fuzzCreate(t *testing.T, name, desc string) (*entity.GameInfo, error) {
	game, err := tt.gameService.Create(&dto.CreateGameDTO{
		Name:        name,
		Description: desc,
	})
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return nil, err
	}
	{
		data, _ := json.MarshalIndent(game, "", "	")
		t.Log("create:", string(data))
	}
	return game, nil
}
func (tt *gameTest) fuzzUpdate(t *testing.T, gameID, name, desc string) (*entity.GameInfo, error) {
	game, err := tt.gameService.Update(gameID, &dto.UpdateGameDTO{
		Name:        name,
		Description: desc,
	})
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return nil, err
	}
	{
		data, _ := json.MarshalIndent(game, "", "	")
		t.Log("update:", string(data))
	}
	return game, nil
}
func (tt *gameTest) fuzzDelete(t *testing.T, gameID string) error {
	err := tt.gameService.Delete(gameID)
	if err != nil {
		{
			data, _ := json.MarshalIndent(err, "", "	")
			t.Log(string(data))
		}
		return err
	}
	return nil
}

func FuzzGame(f *testing.F) {
	// Set path for the game test artifacts
	dataPath := os.Getenv("TEST_DATA_PATH")
	if dataPath == "" {
		f.Fatal("TEST_DATA_PATH must be set")
	}
	tt := newGameTest(filepath.Join(dataPath, "game_fuzz_"+uuid.New().String()))

	f.Fuzz(func(t *testing.T, name1, desc1, name2, desc2 string) {
		if utils.NameToID(name1) == "" || utils.NameToID(name2) == "" {
			// skip
			return
		}

		// Empty list
		err := tt.fuzzList(t, 0)
		if err != nil {
			t.Fatal(err)
		}

		// Create game
		game1, err := tt.fuzzCreate(t, name1, desc1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with game
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, game1.ID, name1, desc1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Update game
		game2, err := tt.fuzzUpdate(t, utils.NameToID(name1), name2, desc2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// List with game
		err = tt.fuzzList(t, 1)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Check item
		err = tt.fuzzItem(t, game2.ID, name2, desc2)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Delete game
		err = tt.fuzzDelete(t, utils.NameToID(name2))
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}

		// Empty list
		err = tt.fuzzList(t, 0)
		if err != nil {
			tt.fuzzCleanup() // Cleanup - just in case
			t.Fatal(err)
		}
	})
}