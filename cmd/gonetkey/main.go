package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/app"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus"
	"github.com/skratchdot/open-golang/open"
	"github.com/unathi-skosana/gonetkey/internal/pkg/common"
	"github.com/unathi-skosana/gonetkey/internal/pkg/dialog"
	"github.com/unathi-skosana/gonetkey/internal/pkg/dservice"
	"github.com/unathi-skosana/gonetkey/internal/pkg/icon"
)

func main() {
	app := app.New()

	d := dialog.NewDialog()
	d.Show(app)

	app.Run()

	retries := math.MaxInt64
	inetkey := common.NewInetkey(d.GetUsernameText(), d.GetPasswordText())
	service := dservice.NewDbusService(inetkey)

	service.Run()
	inetkey.Run(retries)
	systray.Run(onReady, onExit)
}

func onExit() {
	fmt.Println("Destroying system stray...")
}

func onReady() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {

		conn, err := dbus.SessionBus()
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		obj := conn.Object(common.BusName, common.ObjectPath)

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

		var status string

		obj.Call(fmt.Sprintf("%s.Status", common.BusName), 0).Store(&status)

		if status != "uninitialized" {
			mOpen.Hide()
		}

		if status == "up" {
			systray.SetTemplateIcon(icon.TrayCitrus, icon.TrayCitrus)
		} else {
			systray.SetTemplateIcon(icon.TrayCitrusRed, icon.TrayCitrusRed)
		}

		toggle := func() {
			if status == "up" {
				obj.Call(fmt.Sprintf("%s.Close", common.BusName), 0)
				status = "down"
				mClose.Hide()
				mOpen.Show()

			} else if status == "down" {
				obj.Call(fmt.Sprintf("%s.Open", common.BusName), 0)
				status = "up"
				mClose.Show()
				mOpen.Hide()
			}
		}

		for {
			select {
			case <-mToggle.ClickedCh:
				toggle()
			case <-mOpen.ClickedCh:
				obj.Call(fmt.Sprintf("%s.Open", common.BusName), 0)
				mClose.Show()
				mOpen.Hide()
			case <-mClose.ClickedCh:
				obj.Call(fmt.Sprintf("%s.Close", common.BusName), 0)
				mClose.Hide()
				mOpen.Show()
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
				obj.Call(fmt.Sprintf("%s.Close", common.BusName), 0)
				systray.Quit()
				return
			}
		}
	}()

	<-c
	systray.Quit()
	return
}
