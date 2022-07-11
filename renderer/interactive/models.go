package interactive

import (
	tea "github.com/charmbracelet/bubbletea"
	apiV1 "github.com/toutdoux-app/toutdoux-cli/api/v1"
)

type model struct {
	todoLists   apiV1.TodoLists
	todoEntries apiV1.TodoEntries
	opts        RendererOptions
}

type RendererOptions struct {
	apiClient           apiV1.Client
	preSelectedTodoList string
}

func NewRendererWithOptions(opts RendererOptions) (*model, error) {

	return &model{
		opts: opts,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}
