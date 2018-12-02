package galo

type PadUID int

type XY struct {
	X float64
	Y float64
}

// Boxer wraps the single Box method.
type Boxer interface {
	// Box returns the bounding box of the given pad.
	Box(paduid PadUID) (bottomLeft, topRight XY)
}
