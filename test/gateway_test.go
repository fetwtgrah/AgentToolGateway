package test

import (
	"encoding/json"
	"errors"
	"testing"

	"AgentToolGateway/server/gateway"
)

// okTool 返回一个会成功执行并返回 "executed" 的工具。
func okTool(name string, risk gateway.RiskLevel) gateway.Tool {
	return gateway.Tool{
		Name: name,
		Risk: risk,
		Execute: func(args json.RawMessage) (any, error) {
			return "executed", nil
		},
	}
}

// newServer 创建一套独立的网关依赖。
func newServer(tool gateway.Tool) (*gateway.GatewayServer, *gateway.ApprovalManager, *gateway.AuditLogger) {
	repo := gateway.NewToolRepo()
	repo.Register(tool)
	am := gateway.NewApprovalManager()
	audit := gateway.NewAuditLogger()
	server := gateway.NewGatewayServer(repo, am, gateway.NewPolicy(), audit)
	return server, am, audit
}

func lastEvent(audit *gateway.AuditLogger) gateway.AuditEvent {
	events := audit.List()
	return events[len(events)-1]
}

func TestLowRiskToolExecutesDirectly(t *testing.T) {
	server, _, audit := newServer(okTool("read_note", gateway.RiskLevelLow))

	result, err := server.HandleToolCall(gateway.ToolCall{
		CallID:   "low-1",
		ToolName: "read_note",
		Args:     json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("低风险工具应直接执行，返回错误: %v", err)
	}
	if result != "executed" {
		t.Fatalf("期望结果 executed，实际为 %v", result)
	}

	ev := lastEvent(audit)
	if ev.Decision != gateway.Allow {
		t.Fatalf("期望决策 allow，实际为 %v", ev.Decision)
	}
	if ev.Status != "success" {
		t.Fatalf("期望状态 success，实际为 %v", ev.Status)
	}
}

func TestMediumRiskToolRequiresApproval(t *testing.T) {
	server, am, audit := newServer(okTool("send_email", gateway.RiskLevelMedium))

	result, err := server.HandleToolCall(gateway.ToolCall{
		CallID:   "medium-1",
		ToolName: "send_email",
		Args:     json.RawMessage(`{"to":"test@qq.com"}`),
	})
	if err != nil {
		t.Fatalf("中风险工具应进入审批而不是报错: %v", err)
	}

	approval, ok := result.(*gateway.Approval)
	if !ok {
		t.Fatalf("期望返回 *gateway.Approval，实际类型为 %T", result)
	}
	if approval.Status != gateway.Pending {
		t.Fatalf("期望审批状态 pending，实际为 %v", approval.Status)
	}

	ev := lastEvent(audit)
	if ev.Decision != gateway.RequiredApproval {
		t.Fatalf("期望决策 required_approval，实际为 %v", ev.Decision)
	}
	if ev.Status != "waiting_approval" {
		t.Fatalf("期望状态 waiting_approval，实际为 %v", ev.Status)
	}

	// 审批通过后执行
	if err := am.Approve(approval.ID); err != nil {
		t.Fatalf("审批通过失败: %v", err)
	}

	executed, err := server.ExecuteApproval(approval.ID)
	if err != nil {
		t.Fatalf("审批通过后执行失败: %v", err)
	}
	if executed != "executed" {
		t.Fatalf("期望执行结果 executed，实际为 %v", executed)
	}

	ev = lastEvent(audit)
	if ev.Status != "success" {
		t.Fatalf("期望状态 success，实际为 %v", ev.Status)
	}
}

func TestHighRiskToolRejected(t *testing.T) {
	server, _, audit := newServer(okTool("delete_note", gateway.RiskLevelHigh))

	_, err := server.HandleToolCall(gateway.ToolCall{
		CallID:   "high-1",
		ToolName: "delete_note",
		Args:     json.RawMessage(`{"path":"draft.txt"}`),
	})
	if err == nil {
		t.Fatal("高风险工具应被拒绝，但未返回错误")
	}
	if err.Error() != "permission denied" {
		t.Fatalf("期望错误 permission denied，实际为 %v", err)
	}

	ev := lastEvent(audit)
	if ev.Decision != gateway.Deny {
		t.Fatalf("期望决策 deny，实际为 %v", ev.Decision)
	}
	if ev.Status != "denied" {
		t.Fatalf("期望状态 denied，实际为 %v", ev.Status)
	}
}

func TestAuditStatusChangesOnExecutionFailure(t *testing.T) {
	failTool := gateway.Tool{
		Name: "fail_tool",
		Risk: gateway.RiskLevelLow,
		Execute: func(args json.RawMessage) (any, error) {
			return nil, errors.New("boom")
		},
	}
	server, _, audit := newServer(failTool)

	_, err := server.HandleToolCall(gateway.ToolCall{
		CallID:   "fail-1",
		ToolName: "fail_tool",
		Args:     json.RawMessage(`{}`),
	})
	if err == nil {
		t.Fatal("工具执行失败时应返回错误")
	}

	ev := lastEvent(audit)
	if ev.Decision != gateway.Allow {
		t.Fatalf("期望决策 allow，实际为 %v", ev.Decision)
	}
	if ev.Status != "fail" {
		t.Fatalf("期望状态 fail，实际为 %v", ev.Status)
	}
}
