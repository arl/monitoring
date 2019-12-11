package main

import (
	"image"
	gopalette "image/color/palette"
	"testing"
)

func BenchmarkRenderFrame(b *testing.B) {
	m := defaultCfg
	img := image.NewPaletted(image.Rect(0, 0, m.width, m.height), gopalette.Plan9)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		m.renderFrame(m.bounds, img)
	}
}
