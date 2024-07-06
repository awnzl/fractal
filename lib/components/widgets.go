package components

import (
	"github.com/awnzl/fractal/lib/fractal"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type FractalRender struct {
	widget.BaseWidget
	fractal fractal.Fractal
}

func NewFractalRender(f fractal.Fractal) *FractalRender {
	render := &FractalRender{
		fractal: f,
	}

	render.ExtendBaseWidget(render)
	return render
}

func (render *FractalRender) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRaster(render.fractal.ImageRender))
}

// Handle scroll events
func (render *FractalRender) Scrolled(ev *fyne.ScrollEvent) {
	render.fractal.Zoom(render.Size().Width, render.Size().Height, ev.Position.X, ev.Position.Y, ev.Scrolled.DY)
	render.Refresh()
}
