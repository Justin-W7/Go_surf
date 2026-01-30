package models

type BearingDegree int

const (
	N  BearingDegree = 0
	NE BearingDegree = 45
	E  BearingDegree = 90
	SE BearingDegree = 135
	S  BearingDegree = 180
	SW BearingDegree = 225
	W  BearingDegree = 270
	NW BearingDegree = 315
)

var BearingMap = map[string]BearingDegree{
	"N":  N,
	"NE": NE,
	"E":  E,
	"SE": SE,
	"S":  S,
	"SW": SW,
	"W":  W,
	"NW": NW,
}
