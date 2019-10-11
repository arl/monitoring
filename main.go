package main

import (
	"image"
	"image/png"
	"log"
	"net/http"
)

type server struct {
	mandelbrot
}

func (s *server) renderFractal() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "image/png")
		img := image.NewRGBA(image.Rect(0, 0, 1920, 1080))
		s.mandelbrot.render(img)
		png.Encode(w, img)
	}
}

func main() {
	s := &server{
		mandelbrot{
			maxIter: 200,
			minReal: -2,
			maxReal: 1,
			minImag: -1,
			maxImag: 1,
		},
	}

	http.Handle("/", s.renderFractal())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
