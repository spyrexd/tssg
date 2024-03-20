package main

import (
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/spyrexd/tssg/internal/components"
	"github.com/spyrexd/tssg/internal/config"
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

	rootPath := "public"
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		log.Fatalf("failed to create output directory %v", err)
	}

	indexPath := path.Join(rootPath, "index.html")
	indexPage, err := os.OpenFile(indexPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("faild to create file: %v", err)
	}

	index := components.Index{Title: "Adveture 1021", Lists: lists}
	err = index.Render(indexPage)
	if err != nil {
		log.Fatalf("faild to render file: %v", err)
	}

}
