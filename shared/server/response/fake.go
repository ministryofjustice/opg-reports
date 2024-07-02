package response

import "opg-reports/shared/fake"

func FakeTable(bodyRows int, cellHeadersCount int, cellDataCount int, cellExtrasCount int) (tb *Table[*Cell, *Row[*Cell]]) {
	head := FakeRow(cellHeadersCount, cellDataCount, cellExtrasCount)
	foot := FakeRow(cellHeadersCount, cellDataCount, cellExtrasCount)
	body := FakeRows(bodyRows, cellHeadersCount, cellDataCount, cellExtrasCount)
	tb = NewTable[*Cell, *Row[*Cell]]()
	tb.SetTableHead(head)
	tb.SetTableFoot(foot)
	tb.SetTableBody(body...)
	return
}

func FakeRows(rows int, cellHeadersCount int, cellDataCount int, cellExtrasCount int) (list []*Row[*Cell]) {
	for i := 0; i < rows; i++ {
		r := FakeRow(cellHeadersCount, cellDataCount, cellExtrasCount)
		list = append(list, r)
	}
	return
}

func FakeRow(cellHeadersCount int, cellDataCount int, cellExtrasCount int) (row *Row[*Cell]) {
	cells := FakeCells(cellHeadersCount, cellDataCount, cellExtrasCount)
	row = NewRow(cells...)
	return
}

func FakeCells(headerCount int, dataCount int, extraCount int) (cells []*Cell) {
	cells = []*Cell{}

	for i := 0; i < headerCount; i++ {
		cells = append(cells, FakeCellHeader())
	}
	for i := 0; i < dataCount; i++ {
		cells = append(cells, FakeCell())
	}
	for i := 0; i < extraCount; i++ {
		cells = append(cells, FakeCellExtra())
	}
	return
}

func FakeCell() *Cell {
	return NewCell(fake.String(5), fake.String(30))
}
func FakeCellHeader() *Cell {
	return NewCellHeader(fake.String(5), fake.String(30))
}
func FakeCellExtra() *Cell {
	return NewCellExtra(fake.String(5), fake.String(30))
}
