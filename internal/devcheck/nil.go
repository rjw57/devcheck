package devcheck

type NilCheck struct{}

func NewNilCheck() *NilCheck {
	return &NilCheck{}
}

func (c *NilCheck) Check(l *Logger) error {
	l.Success("Nil check")
	return nil
}
