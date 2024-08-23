package emit

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"ytdlp/helpers/logrus"
)

type DownloadStatus string

const (
	DownloadStatusPending    DownloadStatus = "pending"
	DownloadStatusProcessing DownloadStatus = "processing"
	DownloadStatusDownload   DownloadStatus = "download"
	DownloadStatusDone       DownloadStatus = "done"
	DownloadStatusError      DownloadStatus = "error"
)

const (
	DownloadStart    = "download-start"
	DownloadStop     = "download-stop"
	DownloadProgress = "download-progress"
)

type EmitDownload struct {
	ctx *context.Context
}

type JsonDownloadStruct struct {
	Status   DownloadStatus `json:"status,omitempty"`
	Message  string         `json:"message,omitempty"`
	Progress interface{}    `json:"progress,omitempty"`
}

func NewEmitDownload(ctx *context.Context) EmitDownload {
	return EmitDownload{
		ctx: ctx,
	}
}

func (e *EmitDownload) Start() {
	emitKey := DownloadStart
	runtime.EventsEmit(*e.ctx, emitKey, JsonDownloadStruct{})
}

func (e *EmitDownload) Stop() {
	emitKey := DownloadStop
	runtime.EventsEmit(*e.ctx, emitKey)
}

func (e *EmitDownload) Progress(status DownloadStatus, progress interface{}) {
	emitKey := DownloadProgress
	logrus.LogrusLoggerWithContext(e.ctx).Debugf("Emitting progress: %s", emitKey)
	runtime.EventsEmit(*e.ctx, emitKey, JsonDownloadStruct{
		Status:   status,
		Progress: progress,
	})
}
