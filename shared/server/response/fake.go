package response

import "opg-reports/shared/fake"

func FakeTable(bodyRows int, cellHeadersCount int, cellDataCount int, cellExtrasCount int) (tb *Table[ICell, IRow[ICell]]) {
	head := FakeRow(cellHeadersCount, cellDataCount, cellExtrasCount)
	foot := FakeRow(cellHeadersCount, cellDataCount, cellExtrasCount)
	body := FakeRows(bodyRows, cellHeadersCount, cellDataCount, cellExtrasCount)
	tb = NewTable[ICell, IRow[ICell]]()
	tb.SetTableHead(head)
	tb.SetTableFoot(foot)
	tb.SetTableBody(body...)
	return
}

func FakeRows(rows int, cellHeadersCount int, cellDataCount int, cellExtrasCount int) (list []IRow[ICell]) {
	for i := 0; i < rows; i++ {
		r := FakeRow(cellHeadersCount, cellDataCount, cellExtrasCount)
		list = append(list, r)
	}
	return
}

func FakeRow(cellHeadersCount int, cellDataCount int, cellExtrasCount int) (row IRow[ICell]) {
	cells := FakeCells(cellHeadersCount, cellDataCount, cellExtrasCount)
	row = NewRow(cells...)
	return
}

func FakeCells(headerCount int, dataCount int, extraCount int) (cells []ICell) {
	cells = []ICell{}

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

func FakeCell() ICell {
	return NewCell(fake.String(5), fake.String(30))
}
func FakeCellHeader() ICell {
	return NewCellHeader(fake.String(5), fake.String(30))
}
func FakeCellExtra() ICell {
	return NewCellExtra(fake.String(5), fake.String(30))
}
