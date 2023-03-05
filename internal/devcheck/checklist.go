package devcheck

import "errors"

type CheckList struct {
	List []Checkable
}

func NewCheckList(args ...Checkable) *CheckList {
	return &CheckList{List: args}
}

func (c *CheckList) Add(chk Checkable) *CheckList {
	c.List = append(c.List, chk)
	return c
}

func (c *CheckList) Check(l *Logger) error {
	var errs []error
	for _, subc := range c.List {
		e := subc.Check(l)
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	} else {
		return errors.Join(errs...)
	}
}
