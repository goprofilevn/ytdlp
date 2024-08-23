package main

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"ytdlp/helpers/logrus"
	"ytdlp/utils/emit"
	"ytdlp/utils/setup"
)

func (a *App) SetupResources() {
	st := setup.NewSetup(&a.ctx)
	if errInstall := st.Install(); errInstall != nil {
		buttons := []string{"Yes", "No"}
		dialogOpts := runtime.MessageDialogOptions{Title: "Install resource Error", Message: errInstall.Error(), Type: runtime.WarningDialog, Buttons: buttons, DefaultButton: "Ok"}
		if _, err := runtime.MessageDialog(a.ctx, dialogOpts); err != nil {
			logrus.LogrusLoggerWithContext(&a.ctx).Error(err.Error())
		}
		return
	}
	runtime.EventsEmit(a.ctx, emit.ResourceFinish)
}
