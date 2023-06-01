package interfaces

type ZKComparable interface {
	Equals(other ZKComparable) bool
}

type PostgresRuleIterator interface {
	Explode() []any
}
