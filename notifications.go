package webutils

import "context"

// Notifier is responsible for publishing notifications
type Notifier interface {
	Notify(
		ctx context.Context,
		opts NotificationOpts,
	) error
}

// NotificationPriority represents a priority
type NotificationPriority int

// NotificationPriorities
const (
	NotificationPriorityLow NotificationPriority = iota + 1
	NotificationPriorityMedium
	NotificationPriorityHigh
)

func (p NotificationPriority) String() string {
	if p == NotificationPriorityLow {
		return "Low"
	}

	if p == NotificationPriorityMedium {
		return "Medium"
	}

	if p == NotificationPriorityHigh {
		return "High"
	}

	return ""
}

// NotificationOpts represents a notification
type NotificationOpts struct {
	Message  string
	Subject  string
	Priority NotificationPriority
	Error    error
	Metadata map[string]string
}

// Validate NotificationOpts
func (opts *NotificationOpts) Validate() error {
	if opts.Message == "" {
		return ErrNoNotificationMessage
	}

	return nil
}
