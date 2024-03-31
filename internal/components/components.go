package components

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/spyrexd/tssg/internal/trello"
)

type Index struct {
	Lists *[]trello.List
	Title string
}

func (idx *Index) Render(writer io.Writer) error {
	err := indexPage(idx).Render(context.Background(), writer)
	if err != nil {
		return err
	}
	return nil
}

type CardPage struct {
	Title       string
	Card        trello.Card
	Attachments *[]*trello.Attachment
}

func (c *CardPage) Render(writer io.Writer) error {
	descHtml, err := c.Card.DescAsHtml()
	if err != nil {
		return err

	}

	body := unsafe(descHtml)
	err = cardPage(c, body).Render(context.Background(), writer)
	if err != nil {
		return err
	}
	return nil
}

// This function allows raw HTML to be rendered ensure the html is trusted
func unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
