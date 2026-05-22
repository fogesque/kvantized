package trello

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"kvantized-bot/internal/domain/taskboard"
)

type Credentials struct {
	APIKey  string
	Token   string
	BoardID string
}

type BoardClient struct {
	httpClient  *http.Client
	credentials Credentials
}

func NewBoardClient(httpClient *http.Client, credentials Credentials) BoardClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return BoardClient{
		httpClient:  httpClient,
		credentials: credentials,
	}
}

func (c BoardClient) ActiveTaskSummary(ctx context.Context) (taskboard.TaskSummary, error) {
	lists, err := c.openLists(ctx)
	if err != nil {
		return taskboard.TaskSummary{}, err
	}

	cards, err := c.openCards(ctx)
	if err != nil {
		return taskboard.TaskSummary{}, err
	}

	return taskboard.NewTaskSummary(lists, cards), nil
}

func (c BoardClient) openLists(ctx context.Context) ([]taskboard.Column, error) {
	var lists []listResponse
	if err := c.fetchJSON(ctx, "lists", map[string]string{
		"filter": "open",
		"fields": "id,name",
	}, &lists); err != nil {
		return nil, fmt.Errorf("fetch Trello lists: %w", err)
	}

	columns := make([]taskboard.Column, 0, len(lists))
	for _, list := range lists {
		columns = append(columns, taskboard.Column{
			ID:   taskboard.ColumnID(list.ID),
			Name: list.Name,
		})
	}

	return columns, nil
}

func (c BoardClient) openCards(ctx context.Context) ([]taskboard.Card, error) {
	var cards []cardResponse
	if err := c.fetchJSON(ctx, "cards/open", map[string]string{
		"fields": "idList",
	}, &cards); err != nil {
		return nil, fmt.Errorf("fetch Trello cards: %w", err)
	}

	tasks := make([]taskboard.Card, 0, len(cards))
	for _, card := range cards {
		tasks = append(tasks, taskboard.Card{
			ColumnID: taskboard.ColumnID(card.IDList),
		})
	}

	return tasks, nil
}

func (c BoardClient) fetchJSON(ctx context.Context, resource string, params map[string]string, output any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.boardURL(resource, params), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("Trello returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		return err
	}

	return nil
}

func (c BoardClient) boardURL(resource string, params map[string]string) string {
	endpoint := url.URL{
		Scheme: "https",
		Host:   "api.trello.com",
		Path:   fmt.Sprintf("/1/boards/%s/%s", url.PathEscape(c.credentials.BoardID), resource),
	}

	query := endpoint.Query()
	query.Set("key", c.credentials.APIKey)
	query.Set("token", c.credentials.Token)
	for key, value := range params {
		query.Set(key, value)
	}
	endpoint.RawQuery = query.Encode()

	return endpoint.String()
}

type listResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type cardResponse struct {
	IDList string `json:"idList"`
}
