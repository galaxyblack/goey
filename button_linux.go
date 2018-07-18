package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type buttonElement struct {
	Control

	onClick clickSlot
	onFocus focusSlot
	onBlur  blurSlot
}

func (w *Button) mount(parent base.Control) (base.Element, error) {
	control, err := gtk.ButtonNewWithLabel(w.Text)
	if err != nil {
		return nil, err
	}
	control.AddEvents(int(gdk.FOCUS_CHANGE_MASK))

	parent.Handle.Add(control)
	control.SetSensitive(!w.Disabled)
	control.SetCanDefault(true)
	if w.Default {
		control.GrabDefault()
	}

	retval := &buttonElement{Control: Control{&control.Widget}}

	control.Connect("destroy", buttonOnDestroy, retval)
	retval.onClick.Set(&control.Widget, w.OnClick)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func buttonOnDestroy(widget *gtk.Button, mounted *buttonElement) {
	mounted.handle = nil
}

func (w *buttonElement) button() *gtk.Button {
	return (*gtk.Button)(unsafe.Pointer(w.handle))
}

func (w *buttonElement) Props() base.Widget {
	button := w.button()
	text, err := button.GetLabel()
	if err != nil {
		panic("Could not get label: " + err.Error())
	}

	return &Button{
		Text:     text,
		Disabled: !button.GetSensitive(),
		Default:  button.HasDefault(),
		OnClick:  w.onClick.callback,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}
}

func (w *buttonElement) updateProps(data *Button) error {
	button := w.button()
	button.SetLabel(data.Text)
	button.SetSensitive(!data.Disabled)
	if data.Default {
		button.GrabDefault()
	}
	w.onClick.Set(w.handle, data.OnClick)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
