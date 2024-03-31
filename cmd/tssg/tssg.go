package main

import (
	"log"
	"log/slog"

	"github.com/spyrexd/tssg/internal/config"
	"github.com/spyrexd/tssg/internal/render"
	"github.com/spyrexd/tssg/internal/trello"
)

func main() {

	config.LoadEnvConfig()

	client := trello.NewTrelloClient()
	boardId, err := trello.GetBoardIdByName(config.Get("TRELLO_BOARD_NAME").(string), client)
	if err != nil {
		slog.Error("Trello Client Error", "error", err)
	}
	log.Printf("Board id: %s", boardId)

	if err := render.Render(boardId, client); err != nil {
		log.Fatalf("Unable to render site %v", err)
	}
}
