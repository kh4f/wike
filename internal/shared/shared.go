package shared

type POINT struct {
	X int32
	Y int32
}

var (
	ScreenW int16
	ScreenH int16
)

func Ptr[T any](v T) *T {
	return &v
}
