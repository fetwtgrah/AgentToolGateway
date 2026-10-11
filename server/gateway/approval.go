package gateway

type ApprovalStatus string

const (
	Pending  ApprovalStatus = "pending"
	Rejected ApprovalStatus = "rejected"
	Allowed  ApprovalStatus = "allowed"
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
		ID:     "123456",
		Call:   call,
		Status: Pending,
	}
	am.approvals[approval.ID] = approval
	return approval
}

func (am *ApprovalManager) Get(id string) *Approval {
	return am.approvals[id]
}
