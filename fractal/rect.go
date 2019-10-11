package main

// rect is a rectangle defined by 2 points.
//
//    x0              x1
// y0 +---------------+
//    |               |
//    |               |
//    |               |
// y1 +---------------+
type rect struct {
	x0, y0 float64
	x1, y1 float64
}

func (r *rect) width() float64  { return r.x1 - r.x0 }
func (r *rect) height() float64 { return r.y1 - r.y0 }

func (r *rect) center() (x, y float64) {
	return (r.x0 + r.x1) / 2, (r.y0 + r.y1) / 2
}

func (r *rect) translate(x, y float64) {
	r.x0 += x
	r.x1 += x
	r.y0 += y
	r.y1 += y
}

func (r *rect) scale(factor float64) {
	r.x0 *= factor
	r.y0 *= factor
	r.x1 *= factor
	r.y1 *= factor
}

func (r *rect) zoom(x, y, zfactor float64) {
	orgx, orgy := r.center()
	// translate to the origin
	r.translate(-orgx, -orgy)
	// zoom/scale
	r.scale(zfactor)
	// translate back to the origin
	r.translate(x, y)
}
