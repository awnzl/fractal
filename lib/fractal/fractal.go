package fractal

import "image"

type Fractal interface {
	ImageRender(w int, h int) image.Image
	Zoom(w, h, x, y, delta float32) //TODO AW: zoom is the same for any fractal, move it from here
}
