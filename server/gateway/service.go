package gateway

import (
	"fmt"
	"time"
)

type GatewayServer struct {
	Repo  *ToolRepo
	Am    *ApprovalManager
	P     *Policy
	Audit *AuditLogger
}

func NewGatewayServer(repo *ToolRepo, am *ApprovalManager, p *Policy, audit *AuditLogger) *GatewayServer {
	return &GatewayServer{
		Repo:  repo,
		Am:    am,
		P:     p,
		Audit: audit,
	}
}

// HandleToolCall 处理一次toolcall
func (server *GatewayServer) HandleToolCall(call ToolCall) (any, error) {
	tool, ok := server.Repo.Get(call.ToolName)
	if !ok {
		return nil, fmt.Errorf("NO TOOL")
	}
	Decision := server.P.Check(tool)
	event := AuditEvent{
		ID:       call.CallID,
		CallID:   call.CallID,
		ToolName: call.ToolName,
		Decision: Decision,
		Status:   "pending",
		CreateAt: time.Now(),
	}
	server.Audit.Log(event)
	switch Decision {
	case Allow:
		result, err := tool.Execute(call.Args)
		if err != nil {
			server.Audit.UpdateStatus(call.CallID, "fail")
			return nil, err
		}
		server.Audit.UpdateStatus(
			call.CallID,
			"success",
		)
		return result, nil
	case RequiredApproval:
		approval := server.Am.Creat(call)
		server.Audit.UpdateStatus(
			call.CallID,
			"waiting_approval",
		)
		return approval, nil
	case Deny:
		server.Audit.UpdateStatus(
			call.CallID,
			"denied",
		)
		return nil, fmt.Errorf("permission denied")
	}

	return nil, nil
}

func (server *GatewayServer) ExecuteApproval(id string) (any, error) {
	approval := server.Am.Get(id)
	if approval == nil {
		return nil, fmt.Errorf("no approval")
	}
	if approval.Status != Approved {
		return nil, fmt.Errorf("not approved")
	}
	call := approval.Call
	tool, ok := server.Repo.Get(call.ToolName)
	if !ok {
		return nil, fmt.Errorf("NO TOOL")
	}
	result, err := tool.Execute(call.Args)
	if err != nil {
		server.Audit.UpdateStatus(call.CallID, "fail")
		return nil, err
	}
	server.Audit.UpdateStatus(
		call.CallID,
		"success",
	)
	return result, nil
}
