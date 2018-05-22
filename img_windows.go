package goey

import (
	win2 "bitbucket.org/rj/goey/syscall"
	"image"
	"image/draw"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func imageToIcon(prop image.Image) (win.HICON, []uint8, error) {
	// Create a mask for the icon.
	// Currently, we are using a straight white mask, but perhaps this
	// should be a copy of the alpha channel if the source image is
	// RGBA.
	bounds := prop.Bounds()
	imgMask := image.NewGray(prop.Bounds())
	draw.Draw(imgMask, bounds, image.White, image.Point{}, draw.Src)
	hmask, _, err := imageToBitmap(imgMask)
	if err != nil {
		return 0, nil, err
	}

	// Convert the image to a bitmap.
	hbitmap, buffer, err := imageToBitmap(prop)
	if err != nil {
		return 0, nil, err
	}

	// Create the icon
	iconinfo := win.ICONINFO{
		FIcon:    win.TRUE,
		HbmMask:  hmask,
		HbmColor: hbitmap,
	}
	hicon := win.CreateIconIndirect(&iconinfo)
	if hicon == 0 {
		panic("Error in CreateIconIndirect")
	}
	return hicon, buffer, nil
}

func imageToBitmap(prop image.Image) (win.HBITMAP, []uint8, error) {
	if img, ok := prop.(*image.RGBA); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Need to convert RGB to BGR
		for i := 0; i < len(buffer); i += 4 {
			buffer[i+0], buffer[i+2] = buffer[i+2], buffer[i+0]
		}

		// Create the bitmap
		hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 4, 8, unsafe.Pointer(&buffer[0]))
		if hbitmap == 0 {
			panic("Error in CreateBitmap")
		}
		return hbitmap, buffer, nil
	} else if img, ok := prop.(*image.Gray); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Create the bitmap
		hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 1, 8, unsafe.Pointer(&img.Pix[0]))
		if hbitmap == 0 {
			panic("Error in CreateBitmap")
		}
		return hbitmap, buffer, nil
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)
	// Need to convert RGB to BGR
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i+0], img.Pix[i+2] = img.Pix[i+2], img.Pix[i+0]
	}

	// Create the bitmap
	hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 4, 8, unsafe.Pointer(&img.Pix[0]))
	if hbitmap == 0 {
		panic("Error in CreateBitmap")
	}
	return hbitmap, img.Pix, nil
}

func (w *Img) mount(parent NativeWidget) (MountedWidget, error) {
	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(w.Image)
	if err != nil {
		return nil, err
	}

	hwnd := win.CreateWindowEx(0, staticClassName, nil,
		win.WS_CHILD|win.WS_VISIBLE|win.SS_BITMAP|win.SS_LEFT,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}
	win.SendMessage(hwnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	retval := &mountedImg{NativeWidget: NativeWidget{hwnd}, imageData: buffer, width: w.Width, height: w.Height}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedImg struct {
	NativeWidget
	imageData     []uint8
	width, height Length
}

func (w *mountedImg) MeasureWidth() (Length, Length) {
	return w.width, w.width
}

func (w *mountedImg) MeasureHeight(width Length) (Length, Length) {
	return w.height, w.height
}

func (w *mountedImg) SetBounds(bounds image.Rectangle) {
	w.NativeWidget.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *mountedImg) updateProps(data *Img) error {
	w.width, w.height = data.Width, data.Height

	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(data.Image)
	if err != nil {
		return err
	}
	w.imageData = buffer
	win.SendMessage(w.hWnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	return nil
}
