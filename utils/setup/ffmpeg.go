package setup

import (
	"context"
	"os"
	"path/filepath"
	goruntime "runtime"
	"ytdlp/helpers/logrus"
	"ytdlp/utils"
	"ytdlp/utils/emit"
	"ytdlp/utils/zip"
)

const (
	FFMPEG_URL = "https://github.com/BtbN/FFmpeg-Builds/releases/download/autobuild-2024-08-20-13-02/ffmpeg-N-116752-g507c2a5774-win64-gpl.zip"
)

type FFmpeg struct {
	ctx *context.Context
}

func NewFFmpeg(ctx *context.Context) *FFmpeg {
	return &FFmpeg{
		ctx: ctx,
	}
}

func (f *FFmpeg) DownloadFFmpeg(emitResource emit.EmitResource) error {
	filePath := filepath.Join(utils.GetDownloadDir(), "ffmpeg.zip")
	if err := utils.CheckOrDeleteFile(filePath); err != nil {
		return err
	}
	return DownloadFile(f.ctx, emitResource, FFMPEG_URL, filePath)
}

func (f *FFmpeg) Setup() error {
	if goruntime.GOOS != "windows" {
		return nil
	}
	filePath := filepath.Join(utils.GetDownloadDir(), "ffmpeg.zip")
	ffmpegDir := filepath.Join(utils.GetResourceDir(), "ffmpeg")
	ffmpegPath := utils.GetFFmpegPath()
	_, errStat := os.Stat(ffmpegPath)
	if os.IsNotExist(errStat) {
		emitResource := emit.NewEmitResource(f.ctx, "ffmpeg", "FFmpeg")
		emitResource.Start()
		if err := utils.CheckOrCreateDir(ffmpegDir); err != nil {
			return err
		}
		if err := f.DownloadFFmpeg(emitResource); err != nil {
			return err
		}
		emitResource.Progress("Extracting", 0)
		if errExtract := zip.UnzipFile(f.ctx, filePath, ffmpegDir); errExtract != nil {
			return errExtract
		}
		emitResource.Stop()
	} else {
		logrus.LogrusLoggerWithContext(f.ctx).Info("FFmpeg is exist")
	}
	return nil
}
