package components

import (
	"bytes"
	"context"
	"io"

	"github.com/a-h/templ"
	tc "github.com/adlio/trello"
	"github.com/yuin/goldmark"
)

type Renderer interface {
	Render(io.Writer) error
}

type List struct {
	List  *tc.List
	Cards *[]Card
}

func (l *List) Render(writer io.Writer) error {
	if err := listComponent(l).Render(context.Background(), writer); err != nil {
		return err
	}
	return nil
}

type Card struct {
	*tc.Card
}

func (c *Card) descAsHtml() (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(c.Desc), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
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

type CardPage struct {
	Title string
	Card  Card
}

func (c *CardPage) Render(writer io.Writer) error {
	descHtml, err := c.Card.descAsHtml()
	if err != nil {
		return err

	}
	body := Unsafe(descHtml)
	err = cardPage(c, body).Render(context.Background(), writer)
	if err != nil {
		return err
	}
	return nil
}

// This function allows raw HTML to be rendered ensure the html is trusted
func Unsafe(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, html)
		return
	})
}
