package response

import (
	"fmt"
	"testing"
)

func TestSharedServerResponseDataTable(t *testing.T) {
	cells := []*Cell{
		NewCell("name1", "v1"),
		NewCell("name1.5", "v1.5"),
		NewCellHeader("name2", "v2"),
		NewCellExtra("name3", "v3"),
	}
	h := NewRow[*Cell](cells...)
	b := NewRow[*Cell](cells...)
	f := NewRow[*Cell](cells...)

	tb := NewTable[*Cell, *Row[*Cell]]()

	tb.SetTableHead(h)
	gh := tb.GetTableHead()
	if h != gh {
		t.Errorf("table head mismatch")
	}

	tb.SetTableBody(b)
	gb := tb.GetTableBody()
	if b != gb[0] {
		t.Errorf("body mismatch")
	}
	tb.SetTableFoot(f)
	gf := tb.GetTableFoot()
	if f != gf {
		t.Errorf("foot mismatch")
	}

}

func TestSharedServerResponseDataRowRaw(t *testing.T) {
	cells := []*Cell{
		NewCell("name1", "v1"),
		NewCell("name1.5", "v1.5"),
		NewCellHeader("name2", "v2"),
		NewCellExtra("name3", "v3"),
	}
	row := NewRow[*Cell]()

	row.SetRaw(cells...)

	if row.GetHeadersCount() != 1 {
		t.Errorf("header set failed")
	}
	if row.GetHeadersCount() != 1 {
		t.Errorf("extra set failed")
	}
	if row.GetDataCount() != 2 {
		t.Errorf("data set failed")
	}

	if row.GetTotalCellCount() != 4 {
		t.Errorf("total count failed")
	}

}
func TestSharedServerResponseDataRowSupplementary(t *testing.T) {
	cells := []*Cell{
		NewCell("name1", "v1"),
		NewCell("name1.5", "v1.5"),
		NewCellHeader("name2", "v2"),
		NewCellExtra("name3", "v3"),
	}
	row := NewRow[*Cell]()

	row.SetSupplementary(cells...)
	d := row.GetSupplementary()

	if len(d) != 1 || row.GetSupplementaryCount() != 1 {
		t.Errorf("failed to add just supplementary")
	}

}
func TestSharedServerResponseDataRowData(t *testing.T) {
	cells := []*Cell{
		NewCell("name1", "v1"),
		NewCell("name1.5", "v1.5"),
		NewCellHeader("name2", "v2"),
		NewCellExtra("name3", "v3"),
	}
	row := NewRow[*Cell]()

	row.SetData(cells...)
	d := row.GetData()

	if len(d) != 2 || row.GetDataCount() != 2 {
		t.Errorf("failed to add just data")
	}

}
func TestSharedServerResponseDataRowHeaders(t *testing.T) {
	cells := []*Cell{
		NewCell("name1", "v1"),
		NewCellHeader("name2", "v2"),
		NewCellExtra("name3", "v3"),
	}
	row := NewRow[*Cell]()

	row.SetHeaders(cells...)
	h := row.GetHeaders()

	if len(h) != 1 || row.GetHeadersCount() != 1 {
		t.Errorf("failed to add just headers")
	}

}
func TestSharedServerResponseDataCellSupplementary(t *testing.T) {

	c1 := NewCell("name1", "v1")
	c2 := NewCellHeader("name2", "v2")
	c3 := NewCellExtra("name3", "v3")

	if c1.GetIsSupplementary() != false {
		t.Errorf("cell is not sup")
	}
	c1.SetIsSupplementary(true)
	if c1.GetIsSupplementary() != true {
		t.Errorf("cell is sup")
	}

	if c2.GetIsSupplementary() != false {
		t.Errorf("cell is not sup")
	}

	if c3.GetIsSupplementary() != true {
		t.Errorf("cell is sup")
	}

}
func TestSharedServerResponseDataCellHeader(t *testing.T) {

	c1 := NewCell("name1", "v1")
	c2 := NewCellHeader("name2", "v2")
	c3 := NewCellExtra("name3", "v3")

	if c1.GetIsHeader() != false {
		t.Errorf("cell is not header")
	}
	c1.SetIsHeader(true)
	if c1.GetIsHeader() != true {
		t.Errorf("cell is header")
	}

	if c2.GetIsHeader() != true {
		t.Errorf("cell is header")
	}

	if c3.GetIsHeader() != false {
		t.Errorf("is not header")
	}

}
func TestSharedServerResponseDataCellValues(t *testing.T) {
	cells := []*Cell{
		NewCell("name1", "v1"),
		NewCellHeader("name2", "v2"),
		NewCellExtra("name3", "v3"),
	}
	for i, c := range cells {
		if c.Value != c.GetValue() {
			t.Errorf("failed to get value")
			fmt.Printf("%+v\n", c)
		}
		n := fmt.Sprintf("test-value-%d", i)
		c.SetValue(n)
		if c.Value != n || c.GetValue() != n {
			t.Errorf("set value failed")
		}
	}

}

func TestSharedServerResponseDataCellNames(t *testing.T) {

	cells := []*Cell{
		NewCell("name1", "v1"),
		NewCellHeader("name2", "v2"),
		NewCellExtra("name3", "v3"),
	}

	for i, c := range cells {
		if c.Name != c.GetName() {
			t.Errorf("failed to get name")
			fmt.Printf("%+v\n", c)
		}
		n := fmt.Sprintf("test-name-%d", i)
		c.SetName(n)
		if c.Name != n || c.GetName() != n {
			t.Errorf("set name failed")
		}

	}

}
