package gateway

import "encoding/json"

type ToolCall struct {
	CallID   string          `json:"call_id"`
	ToolName string          `json:"tool_name"`
	Args     json.RawMessage `json:"args"`
}
type Tool struct {
	Name    string
	Execute func(
		args json.RawMessage,
	) (any, error)
}
type ToolRepo struct {
	tools map[string]Tool
}

func NewToolRepo() *ToolRepo {
	return &ToolRepo{
		tools: make(map[string]Tool),
	}
}
func (repo *ToolRepo) Register(tool Tool) {
	repo.tools[tool.Name] = tool
}
func (repo *ToolRepo) Get(name string) (Tool, bool) {
	tool, ok := repo.tools[name]
	return tool, ok
}
