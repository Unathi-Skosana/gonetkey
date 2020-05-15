package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type dialog struct {
	username *widget.Entry
	password *widget.Entry
	button   *widget.Button
	window   fyne.Window
}

func NewDialog() *dialog {
	d := &dialog{}
	return d
}

func (d *dialog) Show(app fyne.App) {
	// Theme
	app.Settings().SetTheme(&myTheme{})

	// Window
	d.window = app.NewWindow("Gonetkey")
	d.window.SetTitle("Gonetkey")
	d.window.Resize(fyne.NewSize(250, 110))
	d.window.CenterOnScreen()
	d.window.SetPadded(true)

	// Inputs
	d.username = widget.NewEntry()
	d.username.SetPlaceHolder("Username")

	d.password = widget.NewPasswordEntry()
	d.password.SetPlaceHolder("Password")

	content := widget.NewVBox(d.username, d.password)
	content.Append(widget.NewButton("Connect", func() {
		app.Quit()
	}))

	d.window.SetContent(content)
	d.window.Show()
}

func (d *dialog) GetUsernameText() string {
	return d.username.Text
}

func (d *dialog) GetPasswordText() string {
	return d.password.Text
}
