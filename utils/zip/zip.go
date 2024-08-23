package zip

import (
	"archive/zip"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"ytdlp/helpers/logrus"
)

func CompressFolder(ctx *context.Context, inputDir, outputFile string) error {
	outFile, errOutFile := os.Create(outputFile)
	if errOutFile != nil {
		return errOutFile
	}

	w := zip.NewWriter(outFile)

	if err := addFilesToZip(ctx, w, inputDir, ""); err != nil {
		_ = outFile.Close()
		return err
	}

	if err := w.Close(); err != nil {
		_ = outFile.Close()
		return errors.New("Warning: closing zipfile writer failed: " + err.Error())
	}

	if err := outFile.Close(); err != nil {
		return errors.New("Warning: closing zipfile failed: " + err.Error())
	}

	return nil
}

func addFilesToZip(ctx *context.Context, w *zip.Writer, basePath, baseInZip string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		logrus.LogrusLoggerWithContext(ctx).Errorf("Error reading directory %s: %s", basePath, err.Error())
		return err
	}

	for _, file := range files {
		fullfilepath := filepath.Join(basePath, file.Name())
		if _, errStat := os.Stat(fullfilepath); os.IsNotExist(errStat) {
			// ensure the file exists. For example a symlink pointing to a non-existing location might be listed but not actually exist
			continue
		}

		if file.Mode()&os.ModeSymlink != 0 {
			// ignore symlinks alltogether
			continue
		}

		if file.IsDir() {
			if errAdd := addFilesToZip(ctx, w, fullfilepath, filepath.Join(baseInZip, file.Name())); errAdd != nil {
				return errAdd
			}
		} else if file.Mode().IsRegular() {
			dat, errRead := os.ReadFile(fullfilepath)
			if errRead != nil {
				return errRead
			}
			f, errCreate := w.Create(filepath.Join(baseInZip, file.Name()))
			if errCreate != nil {
				return errCreate
			}
			_, errWrite := f.Write(dat)
			if errWrite != nil {
				return errWrite
			}
		} else {
			// we ignore non-regular files because they are scary
		}
	}
	return nil
}
