package main

import (
	"github.com/awnzl/fractal/lib/components"
	"github.com/awnzl/fractal/lib/fractal/mandelbrot"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
)

const (
	initialWidth = 800
	initialHeight = 600
)

func main() {
	a := app.New()
	w := a.NewWindow("Fractals")
	w.SetMaster()
	w.SetFixedSize(false)
	setKeysBindings(w)
	// w.SetMainMenu(makeMenu(a, w))

	r := components.NewFractalRender(mandelbrot.New())
	content := container.NewStack(r)
	w.SetContent(content)
	w.Resize(fyne.NewSize(initialWidth, initialHeight))

	w.ShowAndRun()
}

//TODO AW: add binding to switch between sequential and concurrent logic
func setKeysBindings(w fyne.Window) {
	w.Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyW,
			Modifier: fyne.KeyModifierShortcutDefault,
		},
		func(shortcut fyne.Shortcut) {
			w.Close()
		},
	)
	w.Canvas().AddShortcut(
		&desktop.CustomShortcut{
			KeyName:  fyne.KeyF,
			Modifier: fyne.KeyModifierShortcutDefault,
		},
		func(shortcut fyne.Shortcut) {
			w.SetFullScreen(!w.FullScreen())
		},
	)
	w.Canvas().Capture()
}
