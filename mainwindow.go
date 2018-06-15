package goey

import (
	"errors"
	"image"
	"sync/atomic"
)

var (
	// ErrSetChildrenNotReentrant is returned if a reentrant call to
	// the method SetChildren is called.
	ErrSetChildrenNotReentrant = errors.New("method SetChildren is not reentrant")

	insideSetChildren uintptr
)

// Window represents a top-level window that contain other widgets.
type Window struct {
	windowImpl
}

// NewWindow create a new top-level window for the application.
func NewWindow(title string, children []Widget) (*Window, error) {
	return newWindow(title, children)
}

// Close destroys the window, and releases all associated resources.
func (w *Window) Close() {
	w.close()
}

// Alignment returns the vertical and horiztonal alignment properties of
// the window.
func (w *Window) Alignment() (MainAxisAlign, CrossAxisAlign) {
	return w.getAlignment()
}

// Children returns the mounted children for the window.  In general, this
// method should not be used.
func (w *Window) Children() []Element {
	return w.getChildren()
}

// Message returns a message constructor that can be used to build and then
// show a dialog box with a message.
func (w *Window) Message(text string) *Message {
	ret := NewMessage(text)
	w.message(ret)
	return ret
}

// SetAlignment changes the vertical and horiztonal alignment properties of
// the window.  These properties affect the layout of child widgets.  The
// main axis for alignment is vertical, with the cross axis being horizontal.
func (w *Window) SetAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	return w.setAlignment(main, cross)
}

// SetChildren changes the child windows and widgets of the window.  As
// necessary, GUI widgets will be created or destroyed so that the GUI widgets
// match the widgets described by the parameter children.  The
// position of contained widgets will be updated to match the new layout
// properties.
func (w *Window) SetChildren(children []Widget) error {
	// One source of bugs in widgets is when the fire an event when being
	// updated.  This can lead to reentrant calls to SetChildren, typically
	// with incorrect information since the GUI is in an inconsistent state
	// when the event fires.  In short, this method is not reentrant.
	// The following will block changes to different windows, although
	// that shouldn't be susceptible to the same bugs.  Users in that
	// case should use Do to delay updates to other windows, but it shouldn't
	// happen in practice.
	if !atomic.CompareAndSwapUintptr(&insideSetChildren, 0, 1) {
		return ErrSetChildrenNotReentrant
	}
	defer func() {
		atomic.StoreUintptr(&insideSetChildren, 0)
	}()

	// Defer to the platform-specific code
	return w.setChildren(children)
}

// SetIcon changes the icon associated with the window.
func (w *Window) SetIcon(img image.Image) error {
	return w.setIcon(img)
}

// SetTitle changes the caption in the title bar for the main window.
func (w *Window) SetTitle(title string) error {
	return w.setTitle(title)
}
