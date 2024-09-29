package components

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type Fractal interface {
	ImageRender(w int, h int) image.Image
	Zoom(w, h, x, y, delta float32) //TODO AW: zoom is the same for any fractal, move it from here
}

type FractalRender struct {
	widget.BaseWidget
	fractal Fractal
}

func NewFractalRender(f Fractal) *FractalRender {
	render := &FractalRender{
		fractal: f,
	}

	render.ExtendBaseWidget(render)
	return render
}

func (render *FractalRender) CreateRenderer() fyne.WidgetRenderer {
	raster := canvas.NewRaster(render.fractal.ImageRender)
	raster.SetMinSize(fyne.Size{Width: 320, Height: 240})
	return widget.NewSimpleRenderer(raster)
}

// Handle scroll events
func (render *FractalRender) Scrolled(ev *fyne.ScrollEvent) {
	render.fractal.Zoom(render.Size().Width, render.Size().Height, ev.Position.X, ev.Position.Y, ev.Scrolled.DY)
	render.Refresh()
}
