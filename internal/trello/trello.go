package trello

import (
	"bytes"

	tc "github.com/adlio/trello"
	"github.com/spyrexd/tssg/internal/config"
	"github.com/yuin/goldmark"
)

type List struct {
	List  *tc.List
	Cards *[]Card
}

type Card tc.Card

func (c *Card) DescAsHtml() (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(c.Desc), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type TrelloClient struct {
	client tc.Client
}

func NewTrelloClient() TrelloClient {
	apiKey := config.Get("TRELLO_API_KEY")
	token := config.Get("TRELLO_TOKEN")
	return TrelloClient{
		client: *tc.NewClient(apiKey.(string), token.(string)),
	}
}

func (c *TrelloClient) GetBoardIdByName(boardName string) (string, error) {
	boards, err := c.client.GetMyBoards(tc.Defaults())
	if err != nil {
		return "", err
	}

	for _, board := range boards {
		if board.Name == boardName {
			return board.ID, nil
		}
	}

	return "", nil
}

func (c *TrelloClient) GetBoardLists(boardId string) (*[]List, error) {
	board, err := c.client.GetBoard(boardId, tc.Defaults())
	if err != nil {
		return nil, err
	}

	boardLists, err := board.GetLists(tc.Defaults())
	if err != nil {
		return nil, err
	}

	lists := make([]List, len(boardLists))
	for boardIdx, item := range boardLists {
		listCards, err := item.GetCards(tc.Defaults())
		if err != nil {
			return nil, err
		}

		cards := make([]Card, len(listCards))
		for cardIdx, card := range listCards {
			cards[cardIdx] = Card(*card)
		}

		lists[boardIdx] = List{List: item, Cards: &cards}
	}

	return &lists, nil
}
