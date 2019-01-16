// +build gnustep

package dialog

import (
	"bitbucket.org/rj/goey/cocoa"
)

func (m *OpenFile) show() (string, error) {
	retval := cocoa.OpenPanel(m.parent, m.filename)
	return retval, nil
}

// WithParent sets the parent of the dialog box.
func (m *OpenFile) WithParent(parent *cocoa.Window) *OpenFile {
	m.parent = parent
	return m
}
