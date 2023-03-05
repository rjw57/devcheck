package devcheck

type SectionCheck struct {
	Title    string
	Contents Checkable
}

func NewSectionCheck(title string, contents Checkable) *SectionCheck {
	return &SectionCheck{Title: title, Contents: contents}
}

func (s *SectionCheck) Check(l *Logger) error {
	l.Bullet("%v:", s.Title)
	return s.Contents.Check(l.Indented())
}
