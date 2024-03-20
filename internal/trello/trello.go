package trello

import (
	tc "github.com/adlio/trello"
	"github.com/spyrexd/tssg/internal/components"
	"github.com/spyrexd/tssg/internal/config"
)

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

func (c *TrelloClient) GetBoardLists(boardId string) (*[]components.List, error) {
	board, err := c.client.GetBoard(boardId, tc.Defaults())
	if err != nil {
		return nil, err
	}

	boardLists, err := board.GetLists(tc.Defaults())
	if err != nil {
		return nil, err
	}

	componentLists := make([]components.List, len(boardLists))
	for boardIdx, item := range boardLists {
		cards, err := item.GetCards(tc.Defaults())
		if err != nil {
			return nil, err
		}

		componentCards := make([]components.Card, len(cards))
		for cardIdx, card := range cards {
			componentCards[cardIdx] = components.Card{card}
		}

		componentLists[boardIdx] = components.List{List: item, Cards: &componentCards}
	}

	return &componentLists, nil
}
