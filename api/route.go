package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"tts_deck_build/api/web"
	"tts_deck_build/internal/api"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/repository"
	"tts_deck_build/internal/server"
	"tts_deck_build/internal/service"
)

func GetRoutes() *mux.Router {
	routes := mux.NewRouter().StrictSlash(false)

	web.Init(routes)

	cfg := config.GetConfig()

	// game
	gameRepository := repository.NewGameRepository(cfg)
	gameService := service.NewGameService(gameRepository)
	api.RegisterGameServer(routes, server.NewGameServer(gameService))

	// collection
	collectionRepository := repository.NewCollectionRepository(cfg, gameRepository)
	collectionService := service.NewCollectionService(collectionRepository)
	api.RegisterCollectionServer(routes, server.NewCollectionServer(collectionService))

	// deck
	deckRepository := repository.NewDeckRepository(cfg, collectionRepository)
	deckService := service.NewDeckService(deckRepository)
	api.RegisterDeckServer(routes, server.NewDeckServer(deckService))

	// card
	cardService := service.NewCardService(repository.NewCardRepository(cfg, deckRepository))
	api.RegisterCardServer(routes, server.NewCardServer(cardService))

	// image
	api.RegisterImageServer(routes, server.NewImageServer(gameService, collectionService, deckService, cardService))

	// system
	api.RegisterSystemServer(routes, server.NewSystemServer(cfg))

	// generator
	generatorService := service.NewGeneratorService(cfg, gameService, collectionService, deckService, cardService)
	api.RegisterGeneratorServer(routes, server.NewGeneratorServer(generatorService))

	routes.Use(corsMiddleware)
	return routes
}

// CORS headers
func corsSetupHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// CORS Headers middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsSetupHeaders(w)
		next.ServeHTTP(w, r)
	})
}
