package notifier

type NotificationRouter struct {
	channels map[string]Notifier
	rules    []RoutingRule
}

type Severity uint8

const (
	Info Severity = iota
	Warning
	Fatal
)

type RoutingRule struct {
	EventTypes []string
	Severity   Severity
	Channel    string
	Condition  func(event NotificationEvent) bool
}
