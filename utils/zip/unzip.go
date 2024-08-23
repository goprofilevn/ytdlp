package zip

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"ytdlp/helpers/logrus"
)

func UnzipFile(ctx *context.Context, zipFile, destFolder string) error {
	if runtime.GOOS != "windows" {
		return UnzipCmd(zipFile, destFolder)
	}
	reader, errReader := zip.OpenReader(zipFile)
	if errReader != nil {
		return errReader
	}
	defer func(reader *zip.ReadCloser) {
		if errReaderClose := reader.Close(); errReaderClose != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errReaderClose.Error())
		}
	}(reader)

	for _, file := range reader.File {
		if errUnzip := unzip(ctx, file, destFolder); errUnzip != nil {
			return errors.New(fmt.Sprintf("unzip file error: %s", errUnzip.Error()))
		}
	}

	return nil
}

func unzip(ctx *context.Context, file *zip.File, destFolder string) error {
	filePath := filepath.Join(destFolder, file.Name)
	if strings.HasPrefix(file.Name, "/") {
		return fmt.Errorf("invalid file path: %s", file.Name)
	}
	if file.FileInfo().IsDir() {
		if errMkdir := os.MkdirAll(filePath, os.ModePerm); errMkdir != nil {
			return errMkdir
		}
		return nil
	}
	if errMkdir := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); errMkdir != nil {
		return errMkdir
	}
	outFile, errOutFile := os.Create(filePath)
	if errOutFile != nil {
		return errOutFile
	}
	inFile, errInFile := file.Open()
	if errInFile != nil {
		return errInFile
	}
	defer func(inFile io.ReadCloser) {
		if errInFileClose := inFile.Close(); errInFileClose != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errInFileClose.Error())
		}
	}(inFile)
	if _, errCopy := io.Copy(outFile, inFile); errCopy != nil {
		return errCopy
	}
	defer func(outFile *os.File) {
		if errRemove := outFile.Close(); errRemove != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errRemove.Error())
		}
	}(outFile)
	return nil
}

func UntarFile(ctx *context.Context, tarFile, destFolder string) error {
	file, errFile := os.Open(tarFile)
	if errFile != nil {
		return errors.New(fmt.Sprintf("open file error: %s", errFile.Error()))
	}
	defer func(file *os.File) {
		if errReaderClose := file.Close(); errReaderClose != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errReaderClose.Error())
		}
	}(file)

	gzipReader, errGzipReader := gzip.NewReader(file)
	if errGzipReader != nil {
		return errors.New(fmt.Sprintf("gzip reader error: %s", errGzipReader.Error()))
	}
	defer func(gzipReader *gzip.Reader) {
		if errGzipReaderClose := gzipReader.Close(); errGzipReaderClose != nil {
			logrus.LogrusLoggerWithContext(ctx).Error(errGzipReaderClose.Error())
		}
	}(gzipReader)

	tarReader := tar.NewReader(gzipReader)

	for {
		header, errHeader := tarReader.Next()

		if errHeader == io.EOF {
			break
		}

		if errHeader != nil {
			return errHeader
		}

		if errUntar := untar(ctx, tarReader, header, destFolder); errUntar != nil {
			return errors.New(fmt.Sprintf("untar file error: %s", errUntar.Error()))
		}
	}

	return nil
}

func untar(ctx *context.Context, tarReader *tar.Reader, header *tar.Header, destFolder string) error {
	target := filepath.Join(destFolder, header.Name)
	if header.Typeflag == tar.TypeDir {
		if err := os.MkdirAll(target, 0755); err != nil {
			return err
		}
	} else if header.Typeflag == tar.TypeReg {
		outFile, errOutFile := os.Create(target)
		if errOutFile != nil {
			return errOutFile
		}
		if _, errCopy := io.Copy(outFile, tarReader); errCopy != nil {
			return errCopy
		}
		defer func(outFile *os.File) {
			if errRemove := outFile.Close(); errRemove != nil {
				logrus.LogrusLoggerWithContext(ctx).Error(errRemove.Error())
			}
		}(outFile)
	}
	return nil
}

func UntarCmd(tarFile, destFolder string) error {
	cmd := fmt.Sprintf("tar xzf %s --directory %s", tarFile, destFolder)
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		return err
	}
	return nil
}

func UnzipCmd(zipFile, destFolder string) error {
	cmd := fmt.Sprintf("unzip -o %s -d %s", zipFile, destFolder)
	if err := exec.Command("sh", "-c", cmd).Run(); err != nil {
		if strings.Contains(err.Error(), "exit status 1") {
			return nil
		}
		return err
	}
	return nil
}
