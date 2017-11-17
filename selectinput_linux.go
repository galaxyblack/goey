package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type MountedSelectInput struct {
	NativeWidget

	onChange func(int)
	onFocus  func()
	onBlur   func()
}

func (w *SelectInput) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	for _, v := range w.Items {
		control.AppendText(v)
	}

	retval := &MountedSelectInput{
		NativeWidget: NativeWidget{&control.Widget},
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}

	if w.OnChange != nil {
		control.Connect("changed", selectinput_onChanged, retval)
	}
	control.Connect("destroy", selectinput_onDestroy, retval)
	control.Show()

	return retval, nil
}

func selectinput_onChanged(widget *gtk.ComboBoxText, mounted *MountedSelectInput) {
	mounted.onChange(widget.GetActive())
}

func selectinput_onDestroy(widget *gtk.ComboBoxText, mounted *MountedSelectInput) {
	mounted.handle = nil
}

func (w *MountedSelectInput) UpdateProps(data Widget) error {
	panic("not implemented")
}
