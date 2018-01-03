package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedTextInput struct {
	handle *gtk.Entry

	onChange func(string)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *TextInput) mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetText(w.Value)
	control.SetPlaceholderText(w.Placeholder)

	retval := &mountedTextInput{
		handle:   control,
		onChange: w.OnChange,
	}

	control.Connect("destroy", textinput_onDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "changed", textinput_onChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func textinput_onChanged(widget *gtk.Entry, mounted *mountedTextInput) {
	text, err := widget.GetText()
	if err != nil {
		// TODO:  What is the correct reporting here
		return
	}
	mounted.onChange(text)
}

func textinput_onDestroy(widget *gtk.Entry, mounted *mountedTextInput) {
	mounted.handle = nil
}

func (w *mountedTextInput) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedTextInput) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedTextInput) updateProps(data *TextInput) error {
	w.handle.SetText(data.Value)
	w.handle.SetPlaceholderText(data.Placeholder)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&w.handle.Widget, w.shChange, data.OnChange != nil, "-changed", textinput_onChanged, w)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
