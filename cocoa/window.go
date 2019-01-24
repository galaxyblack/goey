package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "image"
import "unsafe"

// Window is a wrapper for a NSWindow.
type Window struct {
	private int
}

type WindowCallbacks interface {
	OnShouldClose() bool
	OnWillClose()
	OnDidResize()
}

var (
	windowCallbacks = make(map[unsafe.Pointer]WindowCallbacks)
)

func NewWindow(title string, width, height uint) *Window {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.windowNew(ctitle, C.unsigned(width), C.unsigned(height))
	return (*Window)(handle)
}

func (w *Window) Close() {
	C.windowClose(unsafe.Pointer(w))
}

func (w *Window) ContentSize() (int, int) {
	size := C.windowContentSize(unsafe.Pointer(w))
	return int(size.width), int(size.height)
}

func (w *Window) ContentView() *View {
	return (*View)(C.windowContentView(unsafe.Pointer(w)))
}

func (w *Window) MakeFirstResponder(c *Control) {
	C.windowMakeFirstResponder(unsafe.Pointer(w), unsafe.Pointer(c))
}

func (w *Window) SetCallbacks(cb WindowCallbacks) {
	windowCallbacks[unsafe.Pointer(w)] = cb
}

func (w *Window) SetContentSize(width, height int) {
	C.windowSetContentSize(unsafe.Pointer(w), C.int(width), C.int(height))
}

func (w *Window) SetMinSize(width, height int) {
	C.windowSetMinSize(unsafe.Pointer(w), C.int(width), C.int(height))
}

func (w *Window) SetIcon(img image.Image) error {
	nsi, err := imageToNSImage(img)
	if err != nil {
		return err
	}

	C.windowSetIconImage(unsafe.Pointer(w), nsi)
	C.imageClose(nsi)
	return nil
}

func (w *Window) SetScrollVisible(horz, vert bool) {
	C.windowSetScrollVisible(unsafe.Pointer(w), toBool(horz), toBool(vert))
}

func (w *Window) SetTitle(title string) {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.windowSetTitle(unsafe.Pointer(w), ctitle)
}

func (w *Window) Title() string {
	return C.GoString(C.windowTitle(unsafe.Pointer(w)))
}

//export windowShouldClose
func windowShouldClose(handle unsafe.Pointer) bool {
	if cb := windowCallbacks[handle]; cb != nil {
		return cb.OnShouldClose()
	}

	return true
}

//export windowWillClose
func windowWillClose(handle unsafe.Pointer) {
	if cb := windowCallbacks[handle]; cb != nil {
		cb.OnWillClose()
	}
	delete(windowCallbacks, handle)
}

//export windowDidResize
func windowDidResize(handle unsafe.Pointer) {
	if cb := windowCallbacks[handle]; cb != nil {
		cb.OnDidResize()
	}
}