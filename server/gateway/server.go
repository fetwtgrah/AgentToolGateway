package gateway

import "fmt"

type GatewayServer struct {
	Repo *ToolRepo
}

func NewGatewayServer(repo *ToolRepo) *GatewayServer {
	return &GatewayServer{
		Repo: repo,
	}
}
func (server *GatewayServer) Execute(call ToolCall) (any, error) {
	tool, ok := server.Repo.Get(call.ToolName)
	if !ok {
		return nil, fmt.Errorf("NO TOOL")
	}
	result, err := tool.Execute(call.Args)
	if err != nil {
		return nil, err
	}
	return result, nil
}
