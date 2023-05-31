package interfaces

type ZKComparable interface {
	Equals(other ZKComparable) bool
}
