package mandelbrot

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
	"time"
)

var workersNum = runtime.NumCPU() * 2

type ComplexPoint struct {
	r, i float64
	x, y float64
}

type PositionParameters struct {
	kx, ky, centx, centy float64
}

//TODO AW: move parameters to a dedicated package
type Parameters struct {
	depth int
	zoom float64
	rph, gph, bph float64 // Phases for RGB channels
	rfr, gfr, bfr float64 // Frequencies for RGB channels
	kx, ky, centx, centy float64 // position

	position chan PositionParameters // buffered (1) channel for position params
}

func (p *Parameters) Zoom(w, h, x, y, delta float32) {
	var mult float64 = 1.1
	if delta < 0 {
		mult = 0.9
	}

	p.zoom = p.zoom / mult
	mouseX := float64(x / w) * (p.kx * 2) + p.centx - p.kx
	mouseY := float64(y / h) * (p.ky * 2) + p.centy - p.ky
	xMax := (p.centx - p.kx) * mult + mouseX * (1 - mult)
	xMin := (p.centx + p.kx) * mult + mouseX * (1 - mult)
	yMax := (p.centy - p.ky) * mult + mouseY * (1 - mult)
	yMin := (p.centy + p.ky) * mult + mouseY * (1 - mult)
	p.kx = (xMin - xMax) / 2
	p.ky = (yMin - yMax) / 2
	p.centx = (xMin - p.kx)
	p.centy = (yMin - p.ky)

	p.sendPosition()
}

func (p *Parameters) sendPosition() {
	p.position <- PositionParameters{
		kx: p.kx,
		ky: p.ky,
		centx: p.centx,
		centy: p.centy,
	}
}

//TODO AW: add destroy method to close channel before exit
func NewParameters() Parameters {
	p := Parameters{
		depth: 500,
		zoom: 1,
		rph: 4,
		gph: 2,
		bph: 1,
		rfr: 0.15,
		gfr: 0.15,
		bfr: 0.10,
		kx: 1.5,
		ky: 1.2,
		centx: -0.75,
		centy: 0,
		position: make(chan PositionParameters, 1),
	}
	p.sendPosition() // send initial params
	return p
}

type MandelbrotSet struct {
	params Parameters
	name string
	pos PositionParameters
	img *image.RGBA
}

func New() *MandelbrotSet {
	return &MandelbrotSet{
		params: NewParameters(),
		name: "Mandelbrot set",
		img: image.NewRGBA(image.Rect(0, 0, 0, 0)),
	}
}

type Range struct {
	Start int
	End int
}

// splits the image vertically by number
// returns start-end limits for image parts
func (s *MandelbrotSet) GetRanges(partsNum, width int) []Range {
	limits := []Range{}
	partSize := width / partsNum

	for currentPart := 1; currentPart <= partsNum; currentPart++ {
		partEnd := currentPart * partSize
		mult := currentPart - 1
		if mult == 0 {
			limits = append(limits, Range{0, partEnd})
			continue
		}

		partStart := mult * partSize
		if currentPart == partsNum {
			limits = append(limits, Range{partStart, width}) // the last one part
			continue
		}

		limits = append(limits, Range{partStart, partEnd})
	}

	return limits
}

func (s *MandelbrotSet) ImageRender(w int, h int) image.Image {
	if s.img.Rect.Max.X != w || s.img.Rect.Max.Y != h {
		s.img = image.NewRGBA(image.Rect(0, 0, w, h))
	}

	now := time.Now()

	select {
	case s.pos = <-s.params.position:
	default: // window resize case
	}

	//TODO AW: add keys to switch between sequential and concurrent approaches
	//TODO AW: bind few predefined scene positions to some keys to compare rendering time for these two approaches

	//s.sequential(w, h)
	s.concurrent(w, h)

	fmt.Println("time:", time.Since(now))
	return s.img
}

func (s *MandelbrotSet) sequential(w int, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			s.img.SetRGBA(x, y, s.pixel(w, h, x, y))
		}
	}
}

func (s *MandelbrotSet) concurrent(w int, h int) {
	wg := &sync.WaitGroup{}
	wg.Add(workersNum)

	for _, pair := range s.GetRanges(workersNum, w) {
		go func(start, end int) {
			defer wg.Done()
			for y := 0; y < h; y++ {
				for x := start; x < end; x++ {
					s.img.SetRGBA(x, y, s.pixel(w, h, x, y))
				}
			}
		}(pair.Start, pair.End)
	}

	wg.Wait()
}

func (s *MandelbrotSet) Zoom(w, h, x, y, delta float32) {
	s.params.Zoom(w, h, x, y, delta)
}

func (s *MandelbrotSet) pixel(w, h, x, y int) color.RGBA {
	p := ComplexPoint{
		r: 0, // z
		i: 0, // z
		x: float64(x) / float64(w) * (s.pos.kx * 2) + s.pos.centx - s.pos.kx, // c
		y: float64(y) / float64(h) * (s.pos.ky * 2) + s.pos.centy - s.pos.ky, // c
	}

	q := (p.x * p.x - 0.5 * p.x + 0.0625) + p.y * p.y
	if s.checkMandelbrot(q, p) {
		return color.RGBA{0, 0, 0, 255}
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
			return color.RGBA{0, 0, 0, 255}
		} else {
			return s.getColor(s.ci(iter, p))
		}
	}
}

func (s *MandelbrotSet) checkMandelbrot(q float64, p ComplexPoint) bool {
	if (q * (q + (p.x - 0.25)) < (p.y * p.y) / 4) || (p.x * p.x + 2 * p.x + 1 + p.y * p.y < 0.0625) {
		return true
	}
	return false
}

func (s *MandelbrotSet) getColor(n float64) color.RGBA {
	return color.RGBA{
		R: uint8(math.Sin(n * s.params.rfr + s.params.rph) * 127 + 128),
		G: uint8(math.Sin(n * s.params.gfr + s.params.gph) * 127 + 128),
		B: uint8(math.Sin(n * s.params.bfr + s.params.bph) * 127 + 128),
		A: 255,
	}
}

// a continuous iteration count
func (s *MandelbrotSet) ci(iter int, p ComplexPoint) float64 {
	return float64(iter) + 1 - (math.Log((math.Log(math.Sqrt(p.r * p.r + p.i * p.i)) / 2) / math.Log(2)) / math.Log(2))
}
