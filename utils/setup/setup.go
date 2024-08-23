package setup

import (
	"context"
	"os"
	"sync"
	"ytdlp/utils"
)

type Setup struct {
	ctx *context.Context
}

func NewSetup(ctx *context.Context) *Setup {
	return &Setup{
		ctx: ctx,
	}
}

func (s *Setup) Install() error {
	if errInitFolder := s.initFolder(); errInitFolder != nil {
		return errInitFolder
	}
	errChan := make(chan error, 1)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ffmpeg := NewFFmpeg(s.ctx)
		if errSetup := ffmpeg.Setup(); errSetup != nil {
			errChan <- errSetup
		}
	}()
	go func() {
		defer wg.Done()
		ytdlp := NewYtDlp(s.ctx)
		if errSetup := ytdlp.Setup(); errSetup != nil {
			errChan <- errSetup
		}
	}()

	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (s *Setup) initFolder() error {
	// download
	folderDownload := utils.GetDownloadDir()
	_, errStatDownload := os.Stat(folderDownload)
	if os.IsNotExist(errStatDownload) {
		if errMkdir := os.MkdirAll(folderDownload, 0755); errMkdir != nil {
			return errMkdir
		}
	}
	// resource
	_, errStatResource := os.Stat(utils.GetResourceDir())
	if os.IsNotExist(errStatResource) {
		if errMkdir := os.MkdirAll(utils.GetResourceDir(), 0755); errMkdir != nil {
			return errMkdir
		}
	}
	// tempdir
	_, errStatTempDir := os.Stat(utils.GetTempDir())
	if os.IsNotExist(errStatTempDir) {
		if errMkdir := os.MkdirAll(utils.GetTempDir(), 0755); errMkdir != nil {
			return errMkdir
		}
	}
	return nil
}
