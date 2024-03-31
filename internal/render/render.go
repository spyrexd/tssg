package render

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	tc "github.com/adlio/trello"
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

func Render(boardId string, client *tc.Client) error {

	lists, err := trello.GetBoardLists(boardId, client)
	if err != nil {
		log.Fatalf("unable to get board items %v", err)
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

			attachments, err := card.GetAttachments()
			if err != nil {
				log.Fatalf("failed to render page: %v", err)
			}
			if len(*attachments) > 0 {
				imagePath := path.Join(listPath, "images")
				if err := os.MkdirAll(imagePath, 0755); err != nil {
					log.Fatalf("failed to render page: %v", err)
				}
				for _, attachment := range *attachments {
					imageName := attachment.Name
					imageFilePath := path.Join(imagePath, imageName)

					reader, err := trello.GetAttachment(attachment.URL, client)
					if err != nil {
						return err
					}

					out, err := os.Create(imageFilePath)
					if err != nil {
						log.Fatalf("failed to render page: %v", err)
					}
					_, err = io.Copy(out, reader)
					if err != nil {
						log.Fatalf("failed to render page: %v", err)
					}
					out.Close()
				}
			}

			cardPage := components.CardPage{Title: card.Name, Card: card, Attachments: attachments}
			if err := cardPage.Render(page); err != nil {
				log.Fatalf("failed to revder page: %v", err)
			}
		}
	}
	return nil
}
