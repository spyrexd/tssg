package render

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/gosimple/slug"
	"github.com/spyrexd/tssg/internal/components"
	"github.com/spyrexd/tssg/internal/trello"
)

func renderList(list trello.List, writer io.Writer) error {
	if err := components.ListComponent(list).Render(context.Background(), writer); err != nil {
		return err
	}
	return nil
}

func Render(lists *[]trello.List) error {

	rootPath := "public"
	if err := os.MkdirAll(rootPath, 0755); err != nil {
		log.Fatalf("failed to create output directory %v", err)
	}

	indexPath := path.Join(rootPath, "index.html")
	indexPage, err := os.OpenFile(indexPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("faild to create file: %v", err)
	}

	index := components.Index{Lists: lists, Title: "Adveture 1021"}
	err = index.Render(indexPage)
	if err != nil {
		log.Fatalf("faild to render file: %v", err)
	}

	for _, list := range *lists {
		listSlug := slug.Make(list.List.Name)
		listPath := path.Join(rootPath, listSlug)
		if err := os.MkdirAll(listPath, 0755); err != nil {
			log.Fatalf("failed to create output direcroty %v", err)
		}

		listIndexPath := path.Join(listPath, "index.html")
		listIndex, err := os.OpenFile(listIndexPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("faild to create file: %v", err)
		}
		renderList(list, listIndex)

		for _, card := range *list.Cards {
			cardSlug := slug.Make(card.Name)
			cardPath := path.Join(listPath, fmt.Sprintf("%s.html", cardSlug))
			page, err := os.OpenFile(cardPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatalf("faild to create file: %v", err)
			}

			cardPage := components.CardPage{Title: card.Name, Card: card}
			if err := cardPage.Render(page); err != nil {
				log.Fatalf("failed to redner page: %v", err)
			}
		}
	}
	return nil
}
