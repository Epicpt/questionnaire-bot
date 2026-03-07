package entity

type Photo struct {
	Paths []string
	Type  PhotoType
}
type PhotoType int

const (
	Jpeg PhotoType = iota
	Animation
)
