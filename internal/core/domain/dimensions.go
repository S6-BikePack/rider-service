package domain

type Dimensions struct {
	Width  int
	Height int
	Depth  int
}

func NewDimensions(width, height, depth int) Dimensions {
	return Dimensions{width, height, depth}
}
