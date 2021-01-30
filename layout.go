package hexgrid

import (
	"image"
	"math"
)

type Point64 struct {
	X float64
	Y float64
}

func (p *Point64) ToPoint() image.Point {
	return image.Point{X: int(math.Round(p.X)), Y: int(math.Round(p.Y))}
}

type HexLayout struct {
	orientation Orientation
	scale       Point64 // multiplication factor relative to the canonical hexagon, where the points are on a unit circle
	origin      Point64 // center point for hexagon 0,0
}

type Orientation struct {
	f0, f1, f2, f3, b0, b1, b2, b3, startAngle float64
}

var OrientationPointy Orientation = Orientation{math.Sqrt(3.), math.Sqrt(3.) / 2., 0., 3. / 2., math.Sqrt(3.) / 3., -1. / 3., 0., 2. / 3., 0.5}

var OrientationFlat Orientation = Orientation{3. / 2., 0., math.Sqrt(3.) / 2., math.Sqrt(3.), 2. / 3., 0., -1. / 3., math.Sqrt(3.) / 3., 0.}

func NewLayout(gridOrientation Orientation, scaleX, scaleY float64, origin Point64) (HexLayout){
	return HexLayout{orientation: gridOrientation, scale: Point64{scaleX, scaleY}, origin: origin}
}

// HexToPixel returns the center pixel for a given hexagon an a certain Layout
func (l *HexLayout) HexToPoint(h HexCoord) Point64 {

	M := l.orientation
	size := l.scale
	origin := l.origin
	x := (M.f0*float64(h.q) + M.f1*float64(h.r)) * size.X
	y := (M.f2*float64(h.q) + M.f3*float64(h.r)) * size.Y
	return Point64{ X: x + origin.X, Y: y + origin.Y }
}

// PixelToHex returns the corresponding hexagon axial coordinates for a given pixel on a certain Layout
func (l *HexLayout) PointToHex(p Point64) fractionalHex {

	M := l.orientation
	origin := l.origin
	scale := l.scale

	pt := Point64{(p.X - float64(origin.X)) / scale.X, (p.Y - float64(origin.Y)) / scale.Y}
	q := M.b0*pt.X + M.b1*pt.Y
	r := M.b2*pt.X + M.b3*pt.Y
	return fractionalHex{q, r, -q - r}
}

func (l *HexLayout) HexCornerOffset(c int) Point64 {

	M := l.orientation
	angle := 2. * math.Pi * (M.startAngle - float64(c)) / 6.
	return Point64{l.scale.X * math.Cos(angle), l.scale.Y * math.Sin(angle)}
}

// Gets the corners of the hexagon for the given Layout, starting at the E vertex and proceeding in a CCW order
func (l *HexLayout) HexagonCorners(h HexCoord) []Point64 {

	corners := make([]Point64, 0)
	center := l.HexToPoint( h)

	for i := 0; i < 6; i++ {
		offset := l.HexCornerOffset(i)
		corners = append(corners, Point64{center.X + offset.X, center.Y + offset.Y})
	}
	return corners
}
