package endpoint

type IDisplay interface {
	ColumnMap() map[string]string
}

type Display struct {
	columnMap map[string]string
}

func (d *Display) ColumnMap() map[string]string {
	return d.columnMap
}
