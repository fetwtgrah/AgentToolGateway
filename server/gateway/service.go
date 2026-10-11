package gateway

import (
	"fmt"
)

type GatewayServer struct {
	Repo *ToolRepo
	Am   *ApprovalManager
	P    *Policy
}

func NewGatewayServer(repo *ToolRepo, am *ApprovalManager, p *Policy) *GatewayServer {
	return &GatewayServer{
		Repo: repo,
		Am:   am,
		P:    p,
	}
}

// HandleToolCall 处理一次toolcall
func (server *GatewayServer) HandleToolCall(call ToolCall) (any, error) {
	tool, ok := server.Repo.Get(call.ToolName)
	if !ok {
		return nil, fmt.Errorf("NO TOOL")
	}
	switch server.P.Check(tool) {
	case Allow:
		result, _ := tool.Execute(call.Args)
		return result, nil
	case RequiredApproval:
		approval := server.Am.Creat(call)
		return approval, nil
	case Deny:
		return nil, fmt.Errorf("permission denied")
	}

	return nil, nil
}

func (server *GatewayServer) ExecuteApproval(call ToolCall) (any, error) {
	tool, _ := server.Repo.Get(call.ToolName)
	return tool.Execute(call.Args)
}
