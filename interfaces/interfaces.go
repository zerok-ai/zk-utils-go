package interfaces

type ZKComparable interface {
	Equals(other ZKComparable) bool
}

type DbArgs interface {
	GetAllColumns() []any
}
