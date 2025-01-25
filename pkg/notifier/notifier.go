package notifier

import "context"

type NotificationEvent struct {
	Subject string
	Body    string
}

type Notifier interface {
	Notify(ctx context.Context, event NotificationEvent) error
}
