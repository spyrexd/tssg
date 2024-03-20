package components

import (
	"context"
	"io"

	tc "github.com/adlio/trello"
)

type Renderer interface {
	Render(io.Writer) error
}

type List struct {
	*tc.List
}

type Index struct {
	Lists *[]List
	Title string
}

func (idx *Index) Render(writer io.Writer) error {
	err := indexPage(idx).Render(context.Background(), writer)
	if err != nil {
		return err
	}
	return nil
}

func (l *List) Render(witer io.Writer) error {
	return nil
}

type Card struct {
	*tc.Card
}

func (c *Card) Render(witer io.Writer) error {
	return nil
}
