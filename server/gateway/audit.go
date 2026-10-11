package gateway

import "time"

type AuditEvent struct {
	ID       string    `json:"id"`
	CallID   string    `json:"call_id"`
	ToolName string    `json:"tool_name"`
	Decision Decision  `json:"decision"`
	Status   string    `json:"status"`
	Error    string    `json:"error,omitempty"`
	CreateAt time.Time `json:"create_at"`
}
type AuditLogger struct {
	events []AuditEvent
}

func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		events: make([]AuditEvent, 0),
	}
}
func (a *AuditLogger) Log(event AuditEvent) {
	a.events = append(
		a.events,
		event,
	)
}
func (a *AuditLogger) List() []AuditEvent {
	return a.events
}
func (a *AuditLogger) UpdateStatus(
	id string,
	status string,
) {
	for i := range a.events {
		if a.events[i].ID == id {
			a.events[i].Status = status
			return
		}
	}
}
