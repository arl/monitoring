package main

import (
	"image"
	gopalette "image/color/palette"
	"image/gif"
	"io"
	"log"
	"math"

	"gonum.org/v1/plot/palette"
)

// pointA is one of an infinity of interesting points to zoom in.
const pointA = 0.2721950 + 0.00540474i

var defaultCfg = mandelbrot{
	width: 256, height: 256, // image dimension
	nframes: 50, // number of frame in the animated GIF
	bounds: rect{ // 2D-space for the first image
		x0: -2, y0: -1,
		x1: 1, y1: 1,
	},
	zoomLevel: 0.93,   // zoom to apply between a frame and the next one
	zoomPt:    pointA, // 2D coordinates of the complex number to zoom at
	maxiter:   1024,   // number of iteration to check if a pixel is in the mandelbrot set
}

type mandelbrot struct {
	width, height int        // rendered image dimensions
	maxiter       int        // maximum number of iterations
	nframes       int        // how many frames to render
	zoomLevel     float64    // zoom applied at each frame
	zoomPt        complex128 // zoom point

	bounds rect
}

// modulus of a complex number
func mod(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

func (m *mandelbrot) renderFrame(cbounds rect, img *image.Paletted) {
	values := make([]float64, img.Bounds().Dx()*img.Bounds().Dy())
	histogram := make(map[int]float64, img.Bounds().Dx()*img.Bounds().Dy())

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			// Create the complex number corresponding to pixel (x, y)
			c := complex(
				cbounds.x0+(float64(x)*cbounds.width()/float64(img.Bounds().Dx())),
				cbounds.y0+(float64(y)*cbounds.height()/float64(img.Bounds().Dy())),
			)

			const escapeRadius = 2
			const sqEscapeRadius = escapeRadius * escapeRadius

			// Check if c is in the mandelbrot set. In theory we should apply an
			// infinity of iterations. In pratice we know that if z gets bigger
			// than a predefined escape radius, it won't get back and just escape farther.
			var (
				n       int
				z       = 0i
				modulus float64
				v       float64
				escaped bool
			)
			for {
				// modulus = mod(z)
				modulus = real(z)*real(z) + imag(z)*imag(z)
				if modulus >= sqEscapeRadius {
					v = float64(n+1) - math.Log(math.Log2(modulus))
					escaped = true
					break
				}
				if n >= m.maxiter {
					v = float64(m.maxiter)
					break
				}
				z = z*z + c
				n++
			}

			// Record the escape value for that pixel
			values[x+y*img.Bounds().Dy()] = v
			if escaped {
				histogram[int(v)]++
			}
		}
	}

	// Color each pixel
	var total float64
	for _, v := range histogram {
		total += v
	}
	hues := make([]float64, m.maxiter+2)
	var h float64
	i := 0
	for ; i < m.maxiter; i++ {
		h += float64(histogram[i]) / float64(total)
		hues[i] = h
	}
	hues[i], hues[i+1] = h, h

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			mu := values[x+y*img.Bounds().Dy()]
			value := float64(0)
			if mu < float64(m.maxiter) {
				value = 1
			}

			hsva := palette.HSVA{
				H: interpolate(hues[int(math.Floor(mu))], hues[int(math.Ceil(mu))], math.Mod(mu, 1)),
				S: 1,
				V: value,
				A: 1,
			}
			img.Set(x, y, hsva)
		}
	}
}

func interpolate(c1, c2, t float64) float64 {
	return c1*(1-t) + c2*t
}

func (m *mandelbrot) renderAnimatedGif(w io.Writer) {
	images := make([]*image.Paletted, m.nframes)
	delays := make([]int, m.nframes)

	log.Printf("Rendering %d frames", m.nframes)

	bounds := make([]rect, m.nframes)
	bounds[0] = m.bounds
	for i := 1; i < m.nframes; i++ {
		bounds[i] = bounds[i-1]
		bounds[i].zoom(real(m.zoomPt), imag(m.zoomPt), m.zoomLevel)
	}

	for i := 0; i < m.nframes; i++ {
		img := image.NewPaletted(image.Rect(0, 0, m.width, m.height), gopalette.Plan9)
		m.renderFrame(bounds[i], img)
		images[i] = img
	}

	log.Println("Encoding to GIF")

	gif.EncodeAll(w, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}
