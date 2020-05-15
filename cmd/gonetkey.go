package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/app"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"github.com/unathi-skosana/gonetkey/pkg/dialog"
	"github.com/unathi-skosana/gonetkey/pkg/icon"
)

func main() {
	app := app.New()
	app.SetIcon(icon.CitrusBitmap)
	d := dialog.NewDialog()

	d.Show(app)
	app.Run()

	systray.Run(onReady, onExit)
}

func onExit() {
	fmt.Println("Close firewall here")
}

func onReady() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {

		systray.SetTemplateIcon(icon.TrayCitrus, icon.TrayCitrus)
		systray.SetTitle("Goinetkey")
		systray.SetTooltip("Goinetkey")

		mToggle := systray.AddMenuItem("Toggle Firewall", "Toggle Firewall")
		mOpen := systray.AddMenuItem("Open Firewall", "Open Firewall")
		mClose := systray.AddMenuItem("Close Firewall", "Close Firewall")

		systray.AddSeparator()

		mEdit := systray.AddMenuItem("Edit Config", "Edit Config")

		systray.AddSeparator()
		mAdmin := systray.AddMenuItem("User Admin", "User Admin")
		mUsage := systray.AddMenuItem("Usage", "Usage")
		mTariff := systray.AddMenuItem("Tariff", "Tariff")
		mQuit := systray.AddMenuItem("Quit", "Quit")

		shown := true
		toggle := func() {
			if shown {
				systray.SetTemplateIcon(icon.TrayCitrusRed, icon.TrayCitrusRed)
				shown = false
			} else {
				systray.SetTemplateIcon(icon.TrayCitrus, icon.TrayCitrus)
				shown = true
			}
		}

		for {
			select {
			case <-mToggle.ClickedCh:
				toggle()
			case <-mOpen.ClickedCh:
				systray.SetTemplateIcon(icon.TrayCitrus, icon.TrayCitrus)
			case <-mClose.ClickedCh:
				systray.SetTemplateIcon(icon.TrayCitrusRed, icon.TrayCitrusRed)
			case <-mEdit.ClickedCh:
				var home string = os.Getenv("HOME")
				open.Run(fmt.Sprintf("file://%s/.inetkeyrc", home))
			case <-mAdmin.ClickedCh:
				open.Run("https://www.sun.ac.za/useradm")
			case <-mUsage.ClickedCh:
				open.Run("https://maties2.sun.ac.za/fwusage/")
			case <-mTariff.ClickedCh:
				open.Run("https://stbsp01.stb.sun.ac.za/innov/it/it-help/Wiki%20Pages/Internet%20Tariff%20Structure.aspx")
			case <-mQuit.ClickedCh:
				systray.Quit()
				return

			}
		}
	}()

	<-c
	systray.Quit()
	return
}
