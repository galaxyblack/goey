package goey

import (
	"image/color"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	decoration struct {
		className *uint16
		atom      win.ATOM
	}
)

func init() {
	var err error
	decoration.className, err = syscall.UTF16PtrFromString("GoeyBackground")
	if err != nil {
		panic(err)
	}
}

func (w *Decoration) mount(parent base.Control) (base.Element, error) {
	if decoration.atom == 0 {
		var wc win.WNDCLASSEX
		wc.CbSize = uint32(unsafe.Sizeof(wc))
		wc.HInstance = win.GetModuleHandle(nil)
		wc.LpfnWndProc = syscall.NewCallback(decorationWindowProc)
		wc.HCursor = win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW))))
		wc.HbrBackground = win.GetSysColorBrush(win.COLOR_3DFACE)
		wc.LpszClassName = decoration.className

		atom := win.RegisterClassEx(&wc)
		if atom == 0 {
			return nil, syscall.GetLastError()
		}
		decoration.atom = atom
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE)
	hwnd := win.CreateWindowEx(win.WS_EX_CONTROLPARENT, decoration.className, nil, style,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	retval := &decorationElement{
		Control: Control{hwnd},
		fill:    w.Fill,
		stroke:  w.Stroke,
		insets:  w.Insets,
		radius:  w.Radius,
		hBrush:  createBrush(w.Fill),
		hPen:    createPen(w.Stroke),
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	if w.Child != nil {
		child, err := w.Child.Mount(Control{hwnd})
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}
		retval.child = child
	}

	return retval, nil
}

type decorationElement struct {
	Control
	fill   color.RGBA
	stroke color.RGBA
	insets Insets
	radius Length
	hBrush win.HBRUSH
	hPen   win.HPEN

	child     Element
	childSize Size
}

func createBrush(clr color.RGBA) win.HBRUSH {
	// This function create a brush for the requested color.
	//
	// If the color is either white or black, then the stock brush is returned.
	// Note that these can safely be passed to DeleteObject, where they will
	// be ignored.  So we can safely interchange calls to GetStockObject and
	// CreateBrushIndirect.

	if clr.A == 0 {
		// Transparent brush
		return win.HBRUSH(win.GetStockObject(win.NULL_BRUSH))
	} else if clr.R == 0 && clr.G == 0 && clr.B == 0 {
		// Pure black brush
		// TODO:  Implement transparency when clr.A < 0xFF
		return win.HBRUSH(win.GetStockObject(win.BLACK_BRUSH))
	} else if clr.R == 0xff && clr.G == 0xff && clr.B == 0xff {
		// Pure white brush
		// TODO:  Implement transparency when clr.A < 0xFF
		return win.HBRUSH(win.GetStockObject(win.WHITE_BRUSH))
	}

	// There is not stock brush with the correct color.  Create a custom brush.
	lb := win.LOGBRUSH{
		LbStyle: win.BS_SOLID,
		LbColor: win.COLORREF(uint32(clr.B)<<16 | uint32(clr.G)<<8 | uint32(clr.R)),
	}
	return win.CreateBrushIndirect(&lb)
}

func createPen(clr color.RGBA) win.HPEN {
	// This function create a brush for the requested color.
	//
	// If the color is either white or black, then the stock brush is returned.
	// Note that these can safely be passed to DeleteObject, where they will
	// be ignored.  So we can safely interchange calls to GetStockObject and
	// CreateBrushIndirect.

	if clr.A == 0 {
		// Transparent pen
		return win.HPEN(win.GetStockObject(win.NULL_PEN))
	} else if clr.R == 0 && clr.G == 0 && clr.B == 0 {
		// Pure black pen
		// TODO:  Implement transparency when clr.A < 0xFF
		return win.HPEN(win.GetStockObject(win.BLACK_PEN))
	} else if clr.R == 0xff && clr.G == 0xff && clr.B == 0xff {
		// Pure white pen
		// TODO:  Implement transparency when clr.A < 0xFF
		return win.HPEN(win.GetStockObject(win.WHITE_PEN))
	}

	lb := win.LOGBRUSH{
		LbStyle: win.BS_SOLID,
		LbColor: win.COLORREF(uint32(clr.B)<<16 | uint32(clr.G)<<8 | uint32(clr.R)),
	}
	return win.ExtCreatePen(win.PS_COSMETIC|win.PS_SOLID, 1, &lb, 0, nil)
}

func (w *decorationElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
	w.Control.Close()
}

func (w *decorationElement) props() *Decoration {
	// TODO:  Can we determine the color of the brush or pen?  That would allow
	// to verify that the change has propagated right down to the WIN32
	// API.  This code won't detect if their is skew between the colours stored
	// in fill and stroke and the GDI resource hbrush and hpen.

	return &Decoration{
		Fill:   w.fill,
		Stroke: w.stroke,
		Insets: w.insets,
		Radius: w.radius,
	}
}

func (w *decorationElement) SetBounds(bounds base.Rectangle) {
	// Update background control position
	w.Control.SetBounds(bounds)

	if w.child != nil {
		px := FromPixelsX(1)
		py := FromPixelsY(1)
		position := bounds.Min
		bounds.Min.X += px + w.insets.Left - position.X
		bounds.Min.Y += py + w.insets.Top - position.Y
		bounds.Max.X -= px + w.insets.Right + position.X
		bounds.Max.Y -= py + w.insets.Bottom + position.Y
		w.child.SetBounds(bounds)
	}
}

func (w *decorationElement) SetOrder(previous win.HWND) win.HWND {
	previous = w.Control.SetOrder(previous)
	if w.child != nil {
		previous = w.child.SetOrder(previous)
	}
	return previous
}

func (w *decorationElement) updateProps(data *Decoration) error {
	if w.fill != data.Fill {
		// Free the old brush
		if w.hBrush != 0 {
			win.DeleteObject(win.HGDIOBJ(w.hBrush))
		}

		// Allocate the new brush
		w.hBrush = createBrush(data.Fill)
		if w.hBrush == 0 {
			return syscall.GetLastError()
		}
		w.fill = data.Fill
	}

	if w.stroke != data.Stroke {
		if w.hPen != 0 {
			win.DeleteObject(win.HGDIOBJ(w.hPen))
		}

		w.hPen = createPen(data.Stroke)
		if w.hPen == 0 {
			return syscall.GetLastError()
		}
		w.stroke = data.Stroke
	}

	w.insets = data.Insets
	w.radius = data.Radius

	child, err := base.DiffChild(Control{w.hWnd}, w.child, data.Child)
	if err != nil {
		return err
	}
	w.child = child

	return nil
}

func decorationWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		decorationGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_PAINT:
		// Fill with the proper background color
		w := decorationGetPtr(hwnd)

		ps := win.PAINTSTRUCT{}
		cr := win.RECT{}
		win.GetClientRect(hwnd, &cr)
		hdc := win.BeginPaint(hwnd, &ps)
		win.SelectObject(hdc, win.HGDIOBJ(w.hBrush))
		win.SelectObject(hdc, win.HGDIOBJ(w.hPen))
		if w.radius > 0 {
			rx := w.radius.PixelsX()
			ry := w.radius.PixelsY()
			win.RoundRect(hdc, cr.Left, cr.Top, cr.Right, cr.Bottom, int32(rx), int32(ry))
		} else {
			win.Rectangle_(hdc, cr.Left, cr.Top, cr.Right, cr.Bottom)
		}
		win.EndPaint(hwnd, &ps)
		return 0

	case win.WM_COMMAND:
		if n := win.HIWORD(uint32(wParam)); n == win.BN_CLICKED || n == win.EN_UPDATE {
			return win.SendDlgItemMessage(hwnd, int32(win.LOWORD(uint32(wParam))), msg, wParam, lParam)
		}
		// Defer to the default window proc
	}

	// Let the default window proc handle all other messages
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

func decorationGetPtr(hwnd win.HWND) *decorationElement {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*decorationElement)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
