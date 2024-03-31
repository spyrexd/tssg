package trello

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	tc "github.com/adlio/trello"
	"github.com/sirupsen/logrus"
	"github.com/spyrexd/tssg/internal/config"
	"github.com/yuin/goldmark"
)

type List struct {
	List  *tc.List
	Cards *[]Card
}

type Attachment struct {
	*tc.Attachment
}
type Card struct {
	*tc.Card
}

func (c *Card) DescAsHtml() (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(c.Desc), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (c *Card) GetAttachments() (*[]*Attachment, error) {
	attachments, err := c.Card.GetAttachments(tc.Defaults())
	if err != nil {
		return nil, err
	}

	attach := make([]*Attachment, len(attachments))
	for idx, attachment := range attachments {
		attach[idx] = &Attachment{attachment}
	}

	return &attach, nil
}

func NewTrelloClient() *tc.Client {
	apiKey := config.Get("TRELLO_API_KEY")
	token := config.Get("TRELLO_TOKEN")

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	client := *tc.NewClient(apiKey.(string), token.(string))
	return &client

}

func GetBoardIdByName(boardName string, client *tc.Client) (string, error) {
	boards, err := client.GetMyBoards(tc.Defaults())
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

func GetBoardLists(boardId string, client *tc.Client) (*[]List, error) {
	board, err := client.GetBoard(boardId, tc.Defaults())
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
			cards[cardIdx] = Card{Card: card}
		}

		lists[boardIdx] = List{List: item, Cards: &cards}
	}

	return &lists, nil
}

func GetAttachment(downlaodUrl string, client *tc.Client) (io.Reader, error) {

	req, err := http.NewRequest("GET", downlaodUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("OAuth oauth_consumer_key=\"%s\", oauth_token=\"%s\"", client.Key, client.Token))

	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to download attachement, got %v", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil

}
