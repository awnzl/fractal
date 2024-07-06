package mandelbrot

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"
)

type ComplexPoint struct {
	r, i float64
	x, y float64
}

type Parameters struct {
	depth int
	zoom float64
	rph, gph, bph float64 // Phases for RGB channels
	rfr, gfr, bfr float64 // Frequencies for RGB channels
	kx, ky float64
	centx, centy float64
}

type MandelbrotSet struct {
	params Parameters
	name string
}

func NewParameters() Parameters {
	return Parameters{
		depth: 100,
		zoom: 1,
		rph: 4,
		gph: 4,
		bph: 4,
		rfr: 0.15,
		gfr: 0.15,
		bfr: 0.15,
		kx: 1.5,
		ky: 1.2,
		centx: -0.75,
		centy: 0,
	}
}

func New() *MandelbrotSet {
	return &MandelbrotSet{
		params: NewParameters(),
		name: "Mandelbrot set",
	}
}

func (s *MandelbrotSet) ImageRender(w int, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	now := time.Now()

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			s.pixel(w, h, x, y, img)
		}
	}

	fmt.Println("time:", time.Since(now))
	return img
}

//TODO AW: zoom should not be here...
func (s *MandelbrotSet) Zoom(w, h, x, y, delta float32) {
	var mult float64 = 1.1
	if delta < 0 {
		mult = 0.9
	}

	s.params.zoom = s.params.zoom / mult
	mouseX := float64(x / w) * (s.params.kx * 2) + s.params.centx - s.params.kx
	mouseY := float64(y / h) * (s.params.ky * 2) + s.params.centy - s.params.ky
	xMax := (s.params.centx - s.params.kx) * mult + mouseX * (1 - mult)
	xMin := (s.params.centx + s.params.kx) * mult + mouseX * (1 - mult)
	yMax := (s.params.centy - s.params.ky) * mult + mouseY * (1 - mult)
	yMin := (s.params.centy + s.params.ky) * mult + mouseY * (1 - mult)
	s.params.kx = (xMin - xMax) / 2
	s.params.ky = (yMin - yMax) / 2
	s.params.centx = (xMin - s.params.kx)
	s.params.centy = (yMin - s.params.ky)
}

func (s *MandelbrotSet) pixel(w, h, x, y int, img *image.RGBA) { //TODO AW: or it's better to return a value for pixel? check using benchmarks
	p := ComplexPoint{
		r: 0,
		i: 0,
		x: float64(x) / float64(w) * (s.params.kx * 2) + s.params.centx - s.params.kx,
		y: float64(y) / float64(h) * (s.params.ky * 2) + s.params.centy - s.params.ky,
	}

	q := (p.x * p.x - 0.5 * p.x + 0.0625) + p.y * p.y
	if s.checkMandelbrot(q, p) {
		//TODO AW: check if using of img.Pix list is quicker (img.Pix[x + y * size.Width] = 0)
		//TODO AW: since img.Pix is uint8 slice, the formula should account it
		img.Set(x, y, color.Black)
	} else {
		sum := p.r * p.r + p.i * p.i
		sub := p.r * p.r - p.i * p.i
		iter := 0

		for ; sum <= 4 && iter < s.params.depth; {
			p.i = 2 * p.r * p.i + p.y
			p.r = sub + p.x
			sum = p.r * p.r + p.i * p.i
			sub = p.r * p.r - p.i * p.i
			iter++
		}

		if iter == s.params.depth {
			img.Set(x, y, color.Black)
		} else {
			img.Set(x, y, s.getColor(s.ci(iter, p)))
		}
	}
}

func (s *MandelbrotSet) checkMandelbrot(q float64, p ComplexPoint) bool {
	if (q * (q + (p.x - 0.25)) < (p.y * p.y) / 4) || (p.x * p.x + 2 * p.x + 1 + p.y * p.y < 0.0625) {
		return true
	}
	return false
}

func (s *MandelbrotSet) getColor(n float64) color.Color {
	if (n == float64(s.params.depth)) {
		return color.Black
	}
	return color.RGBA{
		R: uint8(math.Sin(n * s.params.rfr + s.params.rph) * 127 + 128),
		G: uint8(math.Sin(n * s.params.gfr + s.params.gph) * 127 + 128),
		B: uint8(math.Sin(n * s.params.bfr + s.params.bph) * 127 + 128),
		A: 255,
	}
}

func (s *MandelbrotSet) ci(iter int, p ComplexPoint) float64 {
	return float64(iter) + 1 - (math.Log((math.Log(math.Sqrt(p.r * p.r + p.i * p.i)) / 2) / math.Log(2)) / math.Log(2))
}
