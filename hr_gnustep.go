// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type hrElement struct {
	control *cocoa.HR
}

func (w *HR) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewHR(parent.Handle)

	retval := &hrElement{
		control: control,
	}
	return retval, nil
}

func (w *hrElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *hrElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}