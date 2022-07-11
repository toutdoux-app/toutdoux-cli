package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type TodoEntry struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TodoListID uuid.UUID `json:"-"`

	Title    string     `json:"title"`
	Priority int        `json:"priority"`
	DueDate  *time.Time `json:"due_date,omitempty"`
	Done     bool       `json:"done"`
	//Labels   TodoListLabels `json:"labels,omitempty"`

	//Relations TodoEntryRelations `many_to_many:"todo_entry_relations" fk_id:"id" db:"-" json:"relations,omitempty"`
}

type TodoEntries []TodoEntry

func (t TodoEntries) Filter(filterFunc func(TodoEntry) bool) TodoEntries {
	return Filter(t, filterFunc)
}

func (t TodoEntries) SplitByDoneStatus() (TodoEntries, TodoEntries) {
	var notDone, done TodoEntries
	for _, todo := range t {
		if todo.Done {
			done = append(done, todo)
		} else {
			notDone = append(notDone, todo)
		}
	}
	return notDone, done
}

const (
	APIGetTodoListEntriesEndpoint = "/list/%s/todo"
)

type GetTodoListEntriesByIDOptions struct{}

func (c *client) GetTodoListEntriesByID(todoListID string, opts GetTodoListEntriesByIDOptions) (TodoEntries, error) {
	var todoEntries TodoEntries

	listListsReq, err := http.NewRequest(http.MethodGet, c.apiURLPathJoin(fmt.Sprintf(APIGetTodoListEntriesEndpoint, todoListID)), nil)
	if err != nil {
		return todoEntries, errors.Wrap(err, "creating http request")
	}

	resp, err := c.httpClient.Do(listListsReq)
	if err != nil {
		return todoEntries, errors.Wrap(err, "performing HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return todoEntries, fmt.Errorf("unexpected HTTP status code %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&todoEntries); err != nil {
		return todoEntries, errors.Wrap(err, "decoding API response")
	}

	return todoEntries, nil
}
