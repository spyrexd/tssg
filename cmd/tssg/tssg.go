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

	trelloClient := trello.NewTrelloClient()
	boardId, err := trelloClient.GetBoardIdByName(config.Get("TRELLO_BOARD_NAME").(string))
	if err != nil {
		slog.Error("Trello Client Error", "error", err)
	}
	log.Printf("Board id: %s", boardId)

	lists, err := trelloClient.GetBoardLists(boardId)
	if err != nil {
		slog.Error("Trello Client Error", "error", err)
	}

	if err := render.Render(lists); err != nil {
		log.Fatalf("Unable to render site %v", err)
	}
}
