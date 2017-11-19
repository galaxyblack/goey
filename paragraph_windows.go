package goey

import (
	"github.com/lxn/win"
	"image"
	"syscall"
	"unsafe"
)

var (
	paragraphMaxMinWidth int
)

func (w *P) Mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16FromString(w.Text)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.SS_LEFT)
	if w.Align == Center {
		style = style | win.SS_CENTER
	} else if w.Align == Right {
		style = style | win.SS_RIGHT
	}

	hwnd := win.CreateWindowEx(0, staticClassName, &text[0],
		style,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	retval := &mountedP{NativeWidget: NativeWidget{hwnd}, text: text}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedP struct {
	NativeWidget
	text []uint16
}

func (w *mountedP) MinimumWidth() DP {
	// If the printed text will be more than 60 characters wide, it will start
	// to impact readability.  We want to force reflow in this case, so we limit
	// the width
	//
	// See the following for the conversion from characters to relative pixels.
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if paragraphMaxMinWidth == 0 {
		hdc := win.GetDC(w.hWnd)
		if hMessageFont != 0 {
			win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
		}
		// Calculate the width of a single 'm' (find the em width)
		rect := win.RECT{0, 0, 0xffff, 0xffff}
		caption := [...]uint16{'m'}
		win.DrawTextEx(hdc, &caption[0], 1, &rect, win.DT_CALCRECT, nil)
		win.ReleaseDC(w.hWnd, hdc)
		paragraphMaxMinWidth = int(rect.Right) * 40
	}

	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, 0xffff, 0xffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT|win.DT_WORDBREAK, nil)
	win.ReleaseDC(w.hWnd, hdc)

	// For reflow if the text is more than 60 characters wide
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if int(rect.Right) > paragraphMaxMinWidth {
		return DP(paragraphMaxMinWidth * 96 / dpi.X)
	}
	return DP(int(rect.Right) * 96 / dpi.X)
}

func (w *mountedP) CalculateHeight(width DP) DP {
	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, int32(width), 0xffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT|win.DT_WORDBREAK, nil)
	win.ReleaseDC(w.hWnd, hdc)
	return DP(int(rect.Bottom) * 96 / dpi.X)
}

func (w *mountedP) SetBounds(bounds image.Rectangle) {
	w.NativeWidget.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *mountedP) UpdateProps(data_ Widget) error {
	data := data_.(*Label)

	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}
	w.text = text
	SetWindowText(w.hWnd, &text[0])
	// TODO:  Update alignment

	return nil
}
