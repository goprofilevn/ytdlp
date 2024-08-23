package setup

import (
	"context"
	"os"
	"path/filepath"
	goruntime "runtime"
	"ytdlp/helpers/logrus"
	"ytdlp/utils"
	"ytdlp/utils/emit"
)

const (
	YT_DLP_URL = "https://github.com/yt-dlp/yt-dlp/releases/download/2024.05.27/yt-dlp.exe"
)

type YtDlp struct {
	ctx *context.Context
}

func NewYtDlp(ctx *context.Context) *YtDlp {
	return &YtDlp{
		ctx: ctx,
	}
}

func (f *YtDlp) DownloadYtDlp(emitResource emit.EmitResource, filePath string) error {
	if err := utils.CheckOrDeleteFile(filePath); err != nil {
		return err
	}
	return DownloadFile(f.ctx, emitResource, YT_DLP_URL, filePath)
}

func (f *YtDlp) Setup() error {
	if goruntime.GOOS != "windows" {
		return nil
	}
	filePath := filepath.Join(utils.GetResourceDir(), "yt-dlp", "yt-dlp.exe")
	youtubeDlDir := filepath.Join(utils.GetResourceDir(), "yt-dlp")
	_, errStat := os.Stat(filePath)
	if os.IsNotExist(errStat) {
		emitResource := emit.NewEmitResource(f.ctx, "yt-dlp", "yt-dlp")
		emitResource.Start()
		if err := utils.CheckOrCreateDir(youtubeDlDir); err != nil {
			return err
		}
		if err := f.DownloadYtDlp(emitResource, filePath); err != nil {
			return err
		}
		emitResource.Stop()
	} else {
		logrus.LogrusLoggerWithContext(f.ctx).Info("yt-dlp already exists")
	}
	return nil
}
