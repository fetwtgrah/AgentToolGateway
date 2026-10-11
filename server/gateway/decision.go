package gateway

type Decision string

const (
	Allow            Decision = "allow"
	RequiredApproval Decision = "required_approval"
	Deny             Decision = "deny"
)

type Policy struct{}

func NewPolicy() *Policy {
	return &Policy{}
}
func (p *Policy) Check(tool Tool) Decision {
	switch tool.Risk {
	case RiskLevelLow:
		return Allow
	case RiskLevelMedium:
		return RequiredApproval
	case RiskLevelHigh:
		return Deny
	default:
		return Deny
	}
}
