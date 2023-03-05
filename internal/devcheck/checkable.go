package devcheck

type Checkable interface {
	Check(*Logger) error
}
