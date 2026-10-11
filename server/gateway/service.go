package gateway

import "fmt"

type GatewayServer struct {
	Repo *ToolRepo
	Am   *ApprovalManager
}

func NewGatewayServer(repo *ToolRepo, am *ApprovalManager) *GatewayServer {
	return &GatewayServer{
		Repo: repo,
		Am:   am,
	}
}

// HandleToolCall 处理一次toolcall
func (server *GatewayServer) HandleToolCall(call ToolCall) (any, error) {
	tool, ok := server.Repo.Get(call.ToolName)
	if !ok {
		return nil, fmt.Errorf("NO TOOL")
	}
	switch tool.Risk {
	case RiskLevelLow:
		result, err := tool.Execute(call.Args)
		if err != nil {
			return nil, err
		}
		return result, nil
	case RiskLevelMedium:
		approval := server.Am.Creat(call)
		return approval, nil

	case RiskLevelHigh:
		return nil, fmt.Errorf("execution denied")
	}
	return nil, nil
}

func (server *GatewayServer) ExecuteApproval(call ToolCall) (any, error) {
	tool, _ := server.Repo.Get(call.ToolName)
	return tool.Execute(call.Args)
}
