package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/gtk"
)

type labelElement struct {
	Control
}

func (w *Label) mount(parent base.Control) (base.Element, error) {
	handle, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	handle.SetSingleLineMode(false)
	parent.Handle.Add(handle)
	handle.SetJustify(gtk.JUSTIFY_LEFT)
	handle.SetHAlign(gtk.ALIGN_START)
	handle.SetLineWrap(false)
	handle.Show()

	retval := &labelElement{Control: Control{&handle.Widget}}
	handle.Connect("destroy", labelOnDestroy, retval)

	return retval, nil
}

func labelOnDestroy(widget *gtk.Label, mounted *labelElement) {
	mounted.handle = nil
}

func (w *labelElement) label() *gtk.Label {
	return (*gtk.Label)(unsafe.Pointer(w.handle))
}

func (w *labelElement) Props() base.Widget {
	label := w.label()
	text, err := label.GetText()
	if err != nil {
		panic("Could not get text, " + err.Error())
	}

	return &Label{
		Text: text,
	}
}

func (w *labelElement) updateProps(data *Label) error {
	label := w.label()
	label.SetText(data.Text)
	return nil
}
