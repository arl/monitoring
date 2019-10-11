package main

import (
	"flag"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	m := defaultCfg

	nframes := flag.Int("frames", m.nframes, "number of frames in final animation")
	dim := flag.Int("dim", m.width, "image dimension")
	zoomLevel := flag.Float64("zoom", m.zoomLevel, "scale to apply at each frame (zoom)")
	maxIter := flag.Int("i", m.maxiter, "max iterations to apply on 𝒛")
	flag.Parse()

	m.nframes = *nframes
	m.width, m.height = *dim, *dim
	m.zoomLevel = *zoomLevel
	m.maxiter = *maxIter

	giff, err := os.Create("out.gif")
	check(err)
	defer giff.Close()

	m.renderAnimatedGif(giff)
}
