package ytdlp

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
	"ytdlp/utils"
	"ytdlp/utils/emit"
)

type SplitState struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type YtDlp struct {
	ctx          *context.Context
	emitDownload emit.EmitDownload
	videoUrl     string
	split        SplitState
}

func NewYtDlp(ctx *context.Context, videoUrl string, split SplitState, emitDownload emit.EmitDownload) *YtDlp {
	return &YtDlp{
		ctx:          ctx,
		videoUrl:     videoUrl,
		split:        split,
		emitDownload: emitDownload,
	}
}

func (y *YtDlp) Download() error {
	select {
	case <-(*y.ctx).Done():
		return fmt.Errorf("context canceled")
	default:
	}
	ytDlpPath := utils.GetYtDlpPath()
	ffmpegPath := utils.GetFFmpegPath()
	cmd := exec.Command(ytDlpPath,
		"--download-sections", fmt.Sprintf("*%v-%v", y.split.Start, y.split.End),
		"--force-keyframes-at-cuts",
		y.videoUrl,
		"--force-overwrites",
		"-S", "res:480,fps",
		"--output", "%UserProfile%\\Downloads\\ytdlp\\%(extractor)s\\%(id)s.%(ext)s",
		"--ffmpeg-location", ffmpegPath,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	_ = cmd.Wait()
	time.Sleep(1 * time.Second)
	return nil
}
