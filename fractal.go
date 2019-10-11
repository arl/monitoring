package main

import (
	"image"
	"math"
)

type mandelbrot struct {
	width, height    int     // rendered image dimensions
	maxIter          int     // maximum number of iterations
	minReal, maxReal float64 // range of real part of c
	minImag, maxImag float64 // range of imaginary part of c
}

func (m *mandelbrot) isBounded(c complex128) int {
	z := complex128(0)
	n := 0
	for abs(z) <= 2 && n < m.maxIter {
		z = z*z + c
		n++
	}
	return n
}

func abs(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

func (m *mandelbrot) render(img *image.RGBA) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			// Convert pixel coordinate to complex number
			c := complex(
				m.minReal+(float64(x)/float64(img.Bounds().Dx()))*(m.maxReal-m.minReal),
				m.minImag+(float64(y)/float64(img.Bounds().Dy()))*(m.maxImag-m.minImag))

			z := complex128(0)
			n := 0
			for abs(z) <= 2 && n < m.maxIter {
				z = z*z + c
				n++
			}

			// Use the number of iterations to scale the Hue component of the color
			col := hsv{h: uint16(65535 * n / m.maxIter), s: 255, v: 0}
			if n < m.maxIter {
				col.v = 255
			}

			img.Set(x, y, col)
		}
	}
}
