package setup

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"ytdlp/helpers/logrus"
	"ytdlp/utils/emit"
)

func getHead(ctx *context.Context, url string) (int64, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		if errClose := Body.Close(); errClose != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errClose.Error())
		}
	}(resp.Body)
	size, errSize := strconv.Atoi(resp.Header.Get("Content-Length"))
	if errSize != nil {
		return 0, errSize
	}
	return int64(size), nil
}

func getProgress(ctx *context.Context, emitResource emit.EmitResource, done chan int64, filePath string, total int64) {
	var stop bool = false
	file, errOpen := os.Open(filePath)
	if errOpen != nil {
		logrus.LogrusLoggerWithContext(ctx).Error(errOpen.Error())
		return
	}
	defer func(File *os.File) {
		if errClose := File.Close(); errClose != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errClose.Error())
		}
	}(file)
	for {
		select {
		case <-done:
			stop = true
		default:
			info, err := file.Stat()
			if err != nil {
				logrus.LogrusLoggerWithContext(ctx).Error(err.Error())
				return
			}
			size := info.Size()
			if size == 0 {
				size = 1
			}
			percent := float64(size) / float64(total) * 100
			emitResource.Progress("Downloading", percent)
		}
		if stop {
			break
		}
		time.Sleep(time.Second)
	}
}

func DownloadFile(ctx *context.Context, emitResource emit.EmitResource, url string, filePath string) error {
	file, fileErr := os.Create(filePath)
	if fileErr != nil {
		return fileErr
	}
	defer func(File *os.File) {
		if err := File.Close(); err != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(err.Error())
		}
	}(file)
	head, errHead := getHead(ctx, url)
	if errHead != nil {
		return errHead
	}
	done := make(chan int64)
	go getProgress(ctx, emitResource, done, filePath, head)
	resp, errResp := http.Get(url)
	if errResp != nil {
		return errResp
	}
	defer func(Body io.ReadCloser) {
		if errClose := Body.Close(); errClose != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errClose.Error())
		}
	}(resp.Body)
	n, errCopy := io.Copy(file, resp.Body)
	if errCopy != nil {
		return errCopy
	}
	done <- n

	return nil
}
