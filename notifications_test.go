package webutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifications(t *testing.T) {
	opts := NotificationOpts{}
	assert.Equal(t, opts.Validate(), ErrNoNotificationMessage)

	opts.Message = "asdf"
	assert.Nil(t, opts.Validate())
}

func TestNotificationPriority(t *testing.T) {
	assert.Equal(t, NotificationPriority(0).String(), "")
	assert.Equal(t, NotificationPriorityLow.String(), "Low")
	assert.Equal(t, NotificationPriorityMedium.String(), "Medium")
	assert.Equal(t, NotificationPriorityHigh.String(), "High")
}
