package gateway

import (
	"fmt"
	"time"
)

type ApprovalStatus string

const (
	Pending  ApprovalStatus = "pending"
	Rejected ApprovalStatus = "rejected"
	Approved ApprovalStatus = "approved"
)

type Approval struct {
	ID     string         `json:"id"`
	Call   ToolCall       `json:"call"`
	Status ApprovalStatus `json:"status"`
}
type ApprovalManager struct {
	approvals map[string]*Approval
}

func NewApprovalManager() *ApprovalManager {
	return &ApprovalManager{
		approvals: make(map[string]*Approval),
	}
}
func (am *ApprovalManager) Creat(call ToolCall) *Approval {
	approval := &Approval{
		ID:     fmt.Sprintf("%v", time.Now().UnixNano()),
		Call:   call,
		Status: Pending,
	}
	am.approvals[approval.ID] = approval
	return approval
}
func (am *ApprovalManager) Get(id string) *Approval {
	return am.approvals[id]
}
func (am *ApprovalManager) Approve(id string) error {
	approval := am.Get(id)
	if approval == nil {
		return fmt.Errorf("approval not found")
	}
	approval.Status = Approved
	return nil
}
func (am *ApprovalManager) Reject(id string) error {
	approval := am.Get(id)
	if approval == nil {
		return fmt.Errorf("approval not found")
	}
	approval.Status = Rejected
	return nil
}
