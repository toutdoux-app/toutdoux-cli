package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type TodoList struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name string `json:"name"`

	UserID uuid.UUID `db:"user_id" json:"user_id"`

	TodoEntries TodoEntries `json:"entries,omitempty"`
	Labels      TodoLabels  `json:"labels,omitempty"`
}

type TodoLists []TodoList

type GetTodoListsOptions struct{}

const (
	APIListTodoListsEndpoint = "/list"
)

func (c *client) GetTodoLists(opts GetTodoListsOptions) (TodoLists, error) {
	var todoLists TodoLists

	listListsReq, err := http.NewRequest(http.MethodGet, c.apiURLPathJoin(APIListTodoListsEndpoint), nil)
	if err != nil {
		return todoLists, errors.Wrap(err, "creating http request")
	}

	resp, err := c.httpClient.Do(listListsReq)
	if err != nil {
		return todoLists, errors.Wrap(err, "performing HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return todoLists, fmt.Errorf("unexpected HTTP status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&todoLists); err != nil {
		return todoLists, errors.Wrap(err, "decoding API response")
	}

	return todoLists, nil
}
