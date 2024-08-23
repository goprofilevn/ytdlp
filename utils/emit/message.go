package emit

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MessageStatus string

const (
	MessageStatusSuccess MessageStatus = "success"
	MessageStatusError   MessageStatus = "error"
	MessageStatusInfo    MessageStatus = "info"
)

type JsonMessageStruct struct {
	Status  MessageStatus `json:"status"`
	Message string        `json:"message"`
}

func Message(ctx *context.Context, status MessageStatus, message string) {
	emitKey := "message"
	runtime.EventsEmit(*ctx, emitKey, JsonMessageStruct{
		Status:  status,
		Message: message,
	})
}
