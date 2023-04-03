package cmp

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type detail struct {
	gtk.Box
	list *gtk.ListBox
}

// creates new detail box base
func newDetail() detail {
	return detail{Box: *gtk.NewBox(gtk.OrientationVertical, 10)}
}

// builds base of the detail box, implements builder pattern
func (d *detail) buildBase(title string) *detail {
	d.SetMarginBottom(20)
	// list holds all detail rows
	d.list = gtk.NewListBox()
	d.list.SetSelectionMode(gtk.SelectionNone)
	d.list.SetShowSeparators(true)
	d.list.AddCSSClass("rich-list")
	// label of the box
	label := gtk.NewLabel(title)
	label.SetMarginBottom(10)
	label.SetHAlign(gtk.AlignStart)
	label.AddCSSClass("title-4")
	// frame is needed to get rounded corners
	frame := gtk.NewFrame("")
	frame.SetChild(d.list)
	// append all
	d.Append(label)
	d.Append(frame)
	return d
}

// add a row to the detail box, implements builder pattern
func (d *detail) addRow(labelText string, value ...gtk.Widgetter) *detail {
	// every row contains a box with label an value
	box := gtk.NewBox(gtk.OrientationHorizontal, 10)
	label := gtk.NewLabel(labelText)
	label.SetHAlign(gtk.AlignStart)
	label.SetHExpand(true)
	box.Append(label)
	for _, w := range value {
		box.Append(w)
	}
	// box is putted into a list row
	row := gtk.NewListBoxRow()
	row.SetChild(box)
	row.SetActivatable(false)
	row.SetSizeRequest(360, 0)
	// append row to list
	d.list.Append(row)
	return d
}
