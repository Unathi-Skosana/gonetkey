package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/app"

	"github.com/getlantern/systray"
	"github.com/godbus/dbus"
	"github.com/unathi-skosana/gonetkey/pkg/common"
	"github.com/unathi-skosana/gonetkey/pkg/dialog"
	"github.com/unathi-skosana/gonetkey/pkg/dservice"
	"github.com/unathi-skosana/gonetkey/pkg/icon"
)

func main() {
	app := app.New()
	app.SetIcon(icon.CitrusBitmap)

	d := dialog.NewDialog()
	d.Show(app)

	app.Run()

	inetkey := common.NewInetkey(d.GetUsernameText(), d.GetPasswordText())
	service := dservice.NewDbusService(inetkey)

	service.Run()
	inetkey.Run(1)
	systray.Run(onReady, onExit)
}

func onExit() {
	fmt.Println("Destroying system stray...")
}

func onReady() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {

		conn, err := dbus.ConnectSessionBus()
		if err != nil {
			panic(err)
		}

		defer conn.Close()

		obj := conn.Object(common.BusName, common.ObjectPath)

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

		var status bool

		obj.Call(fmt.Sprintf("%s.Status", common.BusName), 0).Store(&status)

		mOpen.Hide()

		toggle := func() {
			if status {
				obj.Call(fmt.Sprintf("%s.Close", common.BusName), 0)
				status = false
				mClose.Hide()
				mOpen.Show()

			} else {
				obj.Call(fmt.Sprintf("%s.Open", common.BusName), 0)
				status = true
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
				systray.ShowAppWindow(fmt.Sprintf("file://%s/.inetkeyrc", home))
			case <-mAdmin.ClickedCh:
				systray.ShowAppWindow("https://www.sun.ac.za/useradm")
			case <-mUsage.ClickedCh:
				systray.ShowAppWindow("https://maties2.sun.ac.za/fwusage/")
			case <-mTariff.ClickedCh:
				systray.ShowAppWindow("https://stbsp01.stb.sun.ac.za/innov/it/it-help/Wiki%20Pages/Internet%20Tariff%20Structure.aspx")
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
