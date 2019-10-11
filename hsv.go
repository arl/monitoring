package main

// hsv is a color defined in the HSV plane (Hue-Saturation-Value).
type hsv struct {
	h    uint16
	s, v uint8
}

// RGBA converts an HSV color into RGBA (opaque).
func (c hsv) RGBA() (r, g, b, a uint32) {
	max := uint32(c.v) * 255
	min := uint32(c.v) * uint32(255-c.s)

	c.h %= 360
	segment := c.h / 60
	offset := uint32(c.h % 60)
	mid := ((max - min) * offset) / 60

	switch segment {
	case 0:
		return max, min + mid, min, 0xffff
	case 1:
		return max - mid, max, min, 0xffff
	case 2:
		return min, max, min + mid, 0xffff
	case 3:
		return min, max - mid, max, 0xffff
	case 4:
		return min + mid, min, max, 0xffff
	case 5:
		return max, min, max - mid, 0xffff
	}

	return 0, 0, 0, 0xffff
}
