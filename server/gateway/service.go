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
		return nil,fmt.Errorf("need approval")

	case RiskLevelHigh:
		return nil,fmt.Errorf("execution denied")
	}
	return nil, nil
}
