package dialog

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/unathi-skosana/gonetkey/internal/pkg/icon"
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
	// Icon
	app.SetIcon(icon.CitrusBitmap)

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

	// content box
	content := widget.NewVBox(d.username, d.password)

	// append widgets
	content.Append(widget.NewButton("Connect", func() {
		app.Quit()
	}))

	d.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyReturn:
			app.Quit()
		}
	})

	d.window.SetContent(content)
	d.window.Show()
}

// Get username input text
func (d *dialog) GetUsernameText() string {
	return d.username.Text
}

// Get password input text
func (d *dialog) GetPasswordText() string {
	return d.password.Text
}
