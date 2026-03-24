package config

import (
	"encoding/json"
	"wike/internal/shared"
)

type Region struct {
	X1 int32 `json:"x1"`
	X2 int32 `json:"x2"`
	Y1 int32 `json:"y1"`
	Y2 int32 `json:"y2"`
}

func (r *Region) UnmarshalJSON(data []byte) error {
	type regionAlias Region

	*r = Region{
		X1: 0,
		Y1: 0,
		X2: int32(shared.ScreenWidth),
		Y2: int32(shared.ScreenHeight),
	}

	return json.Unmarshal(data, (*regionAlias)(r))
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
