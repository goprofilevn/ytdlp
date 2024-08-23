package emit

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	ResourceStart    = "resource-start"
	ResourceProgress = "resource-progress"
	ResourceStop     = "resource-stop"
	ResourceError    = "resource-error"
	ResourceFinish   = "resource-finish"
)

type EmitResource struct {
	ctx   *context.Context
	key   string
	title string
}

func NewEmitResource(ctx *context.Context, key string, title string) EmitResource {
	return EmitResource{
		ctx:   ctx,
		key:   key,
		title: title,
	}
}

type JsonResourceStruct struct {
	Key         string  `json:"key"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Progress    float64 `json:"progress"`
}

func (e *EmitResource) Start() {
	emitKey := ResourceStart
	runtime.EventsEmit(*e.ctx, emitKey, JsonResourceStruct{
		Key:         e.key,
		Title:       e.title,
		Description: "",
		Progress:    0,
	})
}

func (e *EmitResource) Progress(description string, progress float64) {
	emitKey := ResourceProgress
	runtime.EventsEmit(*e.ctx, emitKey, JsonResourceStruct{
		Key:         e.key,
		Title:       e.title,
		Description: description,
		Progress:    progress,
	})
}

func (e *EmitResource) Stop() {
	emitKey := ResourceStop
	runtime.EventsEmit(*e.ctx, emitKey, JsonResourceStruct{
		Key:         e.key,
		Title:       e.title,
		Description: "",
		Progress:    100,
	})
}

func (e *EmitResource) Error(message string) {
	emitKey := ResourceError
	runtime.EventsEmit(*e.ctx, emitKey, JsonResourceStruct{
		Key:         e.key,
		Title:       e.title,
		Description: message,
		Progress:    100,
	})
}
