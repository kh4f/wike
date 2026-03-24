package config

import "wike/internal/shared"

type Region struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	W int32 `json:"w"`
	H int32 `json:"h"`
}

func (r *Region) Contains(pt shared.Point) bool {
	rx, ry := r.X, r.Y
	if r.X < 0 {
		rx += int32(shared.ScreenWidth)
	}
	if r.Y < 0 {
		ry += int32(shared.ScreenHeight)
	}

	return pt.X >= rx && pt.X < rx+r.W &&
		pt.Y >= ry && pt.Y < ry+r.H
}
