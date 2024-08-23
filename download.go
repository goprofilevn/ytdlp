package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"ytdlp/helpers/logrus"
	ytdlp "ytdlp/services/yt-dlp"
	"ytdlp/utils/emit"
)

func (a *App) StartDownload(url string, split ytdlp.SplitState) error {
	defer runtime.EventsEmit(a.ctx, emit.DownloadStop)
	ctx, cancel := context.WithCancel(a.ctx)
	a.ctxDownload = ContextState{
		Ctx:    &ctx,
		Cancel: &cancel,
	}
	defer cancel()
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	emitDownload := emit.NewEmitDownload(&ctx)
	ytd := ytdlp.NewYtDlp(&ctx, url, split, emitDownload)
	if err := ytd.Download(); err != nil {
		emit.Message(&ctx, emit.MessageStatusError, err.Error())
		return err
	}
	logrus.LogrusLoggerWithContext(&ctx).Info("Download finished")
	return nil
}
