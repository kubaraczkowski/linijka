package main

import (
	"fmt"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

type Linia struct {
	Index   int
	Bar     string
	checked bool
}

type LiniaModel struct {
	walk.TableModelBase
	items []*Linia
}

func NewLiniaModel() *LiniaModel {
	m := new(LiniaModel)
	m.ExampleRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *LiniaModel) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *LiniaModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index

	case 1:
		return item.Bar

	}

	panic("unexpected col")
}

// Called by the TableView to retrieve if a given row is checked.
func (m *LiniaModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *LiniaModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

func (m *LiniaModel) ExampleRows() {
	m.items = make([]*Linia, 10)

	for i := range m.items {
		m.items[i] = &Linia{
			Index: i,
			Bar:   fmt.Sprintf("Linia %d", i),
		}
	}

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()
}

func main() {

	model := NewLiniaModel()

	var tv *walk.TableView

	MainWindow{
		Title:  "Linijka",
		Size:   Size{800, 600},
		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			TableView{
				AssignTo:         &tv,
				AlternatingRowBG: true,
				CheckBoxes:       true,
				ColumnsOrderable: false,
				MultiSelection:   false,
				Columns: []TableViewColumn{
					{Title: "#"},
					{Title: "Bar"},
				},
				StyleCell: func(style *walk.CellStyle) {
					item := model.items[style.Row()]

					if item.checked {
						if style.Row()%2 == 0 {
							style.BackgroundColor = walk.RGB(159, 215, 255)
						} else {
							style.BackgroundColor = walk.RGB(143, 199, 239)
						}
					}

					switch style.Col() {
					case 1:
						if canvas := style.Canvas(); canvas != nil {
							bounds := style.Bounds()
							bounds.X += 2
							bounds.Y += 2
							bounds.Width = int((float64(bounds.Width) - 4) / 5 * float64(len(item.Bar)))
							bounds.Height -= 4
							bounds.X += 4
							bounds.Y += 2
						}
					}
				},

				Model: model,
				OnSelectedIndexesChanged: func() {
					fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
				},
			},
		},
	}.Run()
}
