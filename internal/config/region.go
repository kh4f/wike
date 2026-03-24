package config

import (
	"wike/internal/shared"

	"gopkg.in/yaml.v3"
)

type Region struct {
	X1 int32 `yaml:"x1"`
	X2 int32 `yaml:"x2"`
	Y1 int32 `yaml:"y1"`
	Y2 int32 `yaml:"y2"`
}

func (r *Region) UnmarshalYAML(value *yaml.Node) error {
	type regionAlias Region

	*r = Region{
		X1: 0,
		Y1: 0,
		X2: int32(shared.ScreenWidth),
		Y2: int32(shared.ScreenHeight),
	}

	return value.Decode((*regionAlias)(r))
}

func (r *Region) Contains(pt shared.Point) bool {
	left, right := orderedBounds(r.X1, r.X2, int32(shared.ScreenWidth))
	top, bottom := orderedBounds(r.Y1, r.Y2, int32(shared.ScreenHeight))

	return pt.X >= left && pt.X < right &&
		pt.Y >= top && pt.Y < bottom
}

func orderedBounds(a int32, b int32, size int32) (int32, int32) {
	a = normalizeCoord(a, size)
	b = normalizeCoord(b, size)
	if a <= b {
		return a, b
	}
	return b, a
}

func normalizeCoord(v int32, size int32) int32 {
	if v < 0 {
		return v + size
	}
	return v
}
