package goey

var (
	columnKind = Kind{"column"}
)

// Column describes a layout widget that arranges its child widgets into a several columns.
// If there is sufficient width, the columns will be arranged side-by-side.  Otherwise, the
// columns will be arranged vertically.
type Column struct {
	Children [][]Widget
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Column) Kind() *Kind {
	return &columnKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Column) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedColumn) Kind() *Kind {
	return &columnKind
}

type mountedColumn struct {
	parent   Control
	children []Element
	counts   []int

	transition Length
}

func (w *Column) mount(parent Control) (Element, error) {
	c := make([]Element, 0, len(w.Children))
	counts := make([]int, 0, len(w.Children))

	for _, v := range w.Children {
		for _, w := range v {
			mountedChild, err := w.Mount(parent)
			if err != nil {
				return nil, err
			}
			c = append(c, mountedChild)
		}
		counts = append(counts, len(v))
	}

	return &mountedColumn{
		parent:   parent,
		children: c,
		counts:   counts,
	}, nil
}

func (w *mountedColumn) Close() {
	// Free the children
	for _, v := range w.children {
		v.Close()
	}
	w.children = nil
}

func (w *mountedColumn) MeasureWidth() (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	gap := calculateHGap(nil, nil)
	min, max := Length(0), Length(0)

	ndx := 0
	for _, v := range w.counts {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children[ndx : ndx+v],
		}
		ndx += v
		tmpMin, tmpMax := vbox.MeasureWidth()

		if tmpMin > min {
			min = tmpMin
		}
		max = max + tmpMax + gap
	}
	w.transition = min.Scale(len(w.counts), 1) + gap.Scale(len(w.counts)-1, 1)
	return min, max
}

func (w *mountedColumn) MeasureHeight(width Length) (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	if w.transition == 0 {
		w.MeasureWidth()
		if w.transition == 0 {
			return 0, 0
		}
	}

	// If now side enough, we will layout the items exactly like a VBox
	if width < w.transition {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children,
		}

		return vbox.MeasureHeight(width)
	}

	ndx := 0
	min, max := Length(0), Length(0)
	for _, v := range w.counts {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children[ndx : ndx+v],
		}
		ndx += v
		tmpMin, tmpMax := vbox.MeasureHeight(width)

		if tmpMin > min {
			min = tmpMin
		}
		if tmpMax > max {
			max = tmpMax
		}
	}
	return min, max
}

func (w *mountedColumn) SetBounds(bounds Rectangle) {
	if len(w.children) == 0 {
		return
	}

	if w.transition == 0 {
		panic("internal error")
	}

	// If not wide enough, we will layout the items exactly like a VBox
	if bounds.Dx() < w.transition {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children,
		}

		vbox.SetBounds(bounds)
		return
	}

	ndx := 0
	count := len(w.counts)
	gap := calculateHGap(nil, nil)
	bounds.Max.X += gap
	for i, v := range w.counts {
		vbox := mountedVBox{
			parent:   w.parent,
			children: w.children[ndx : ndx+v],
		}
		ndx += v

		minx := bounds.Min.X + bounds.Dx().Scale(i, count)
		maxx := bounds.Min.X + bounds.Dx().Scale(i+1, count) - gap
		vbox.SetBounds(Rectangle{Point{minx, bounds.Min.Y}, Point{maxx, bounds.Max.Y}})
	}
}

func (w *mountedColumn) SetChildren(children [][]Widget) error {
	// Flatten list
	c := make([]Widget, 0, len(children))
	for _, v := range children {
		c = append(c, v...)
	}

	err := error(nil)
	w.children, err = DiffChildren(w.parent, w.children, c)
	return err
}

func (w *mountedColumn) updateProps(data *Column) error {
	return w.SetChildren(data.Children)
}

func (w *mountedColumn) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Column))
}
