package goey

import (
	"strconv"
	"testing"

	"bitbucket.org/rj/goey/base"
)

func ExampleButton() {
	clickCount := 0

	// In a full application, this variable would be updated to point to
	// the main window for the application.
	var mainWindow *Window
	// These functions are used to update the GUI.  See below
	var update func()
	var render func() base.Widget

	// Update function
	update = func() {
		err := mainWindow.SetChild(render())
		if err != nil {
			panic(err)
		}
	}

	// Render function generates a tree of Widgets to describe the desired
	// state of the GUI.
	render = func() base.Widget {
		// Prep - text for the button
		text := "Click me!"
		if clickCount > 0 {
			text = text + "  (" + strconv.Itoa(clickCount) + ")"
		}
		// The GUI contains a single widget, this button.
		return &VBox{
			AlignMain:  MainCenter,
			AlignCross: CrossCenter,
			Children: []base.Widget{
				&Button{Text: text, OnClick: func() {
					clickCount++
					update()
				}},
			}}
	}
}

func TestButtonCreate(t *testing.T) {
	testingRenderWidgets(t,
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	)
}

func TestButtonClose(t *testing.T) {
	testingCloseWidgets(t,
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	)
}

func TestButtonFocus(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	)
}

func TestButtonClick(t *testing.T) {
	testingCheckClick(t,
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	)
}

func TestButtonUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	}, []base.Widget{
		&Button{Text: "AB"},
		&Button{Text: "DB", Default: true},
		&Button{Text: "EB", Disabled: true},
	})
}
