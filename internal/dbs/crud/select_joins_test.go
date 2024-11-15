package crud_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ministryofjustice/opg-reports/internal/dbs/adaptors"
	"github.com/ministryofjustice/opg-reports/internal/dbs/crud"
	"github.com/ministryofjustice/opg-reports/internal/fakerextensions/fakerextras"
	"github.com/ministryofjustice/opg-reports/internal/structs"
)

type tJoinBase struct {
	ID int `json:"id,omitempty" db:"id" faker:"-"`
}

func (self *tJoinBase) GetID() int {
	return self.ID
}
func (self *tJoinBase) SetID(id int) {
	self.ID = id
}

type tOrder struct {
	tJoinBase
	ID    int     `json:"id,omitempty" db:"id" faker:"-"`
	Date  string  `json:"date,omitempty" db:"date" faker:"time_string"`
	Items *tItems `json:"items,omitempty" db:"items"`
}

func (self *tOrder) TableName() string {
	return "orders"
}
func (self *tOrder) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "date": "TEXT NOT NULL"}
}
func (self *tOrder) Indexes() map[string][]string {
	return map[string][]string{}
}
func (self *tOrder) InsertColumns() []string {
	return []string{"id", "date"}
}

type tItems []*tItem

func (self *tItems) Scan(src interface{}) (err error) {

	switch src.(type) {
	case []byte:
		err = structs.Unmarshal(src.([]byte), self)
	case string:
		err = structs.Unmarshal([]byte(src.(string)), self)
	default:
		err = fmt.Errorf("unsupported scan src type")
	}
	return
}

type tItem struct {
	tJoinBase
	ID   int    `json:"id,omitempty" db:"id" faker:"-"` // ID is a generated primary key
	Name string `json:"name,omitempty" db:"name" faker:"word"`
}

func (self *tItem) TableName() string {
	return "items"
}
func (self *tItem) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "name": "TEXT NOT NULL"}
}
func (self *tItem) Indexes() map[string][]string {
	return map[string][]string{}
}
func (self *tItem) InsertColumns() []string {
	return []string{"id", "name"}
}

type tOrderItem struct {
	tJoinBase
	ID      int `json:"id,omitempty" db:"id" faker:"-"`
	OrderID int `json:"order_id,omitempty" db:"order_id" faker:"-"`
	ItemID  int `json:"item_id,omitempty" db:"item_id" faker:"-"`
}

func (self *tOrderItem) TableName() string {
	return "order_items"
}
func (self *tOrderItem) Columns() map[string]string {
	return map[string]string{"id": "INTEGER PRIMARY KEY", "order_id": "INTEGER NOT NULL", "item_id": "INTEGER NOT NULL"}
}
func (self *tOrderItem) Indexes() map[string][]string {
	return map[string][]string{}
}
func (self *tOrderItem) InsertColumns() []string {
	return []string{"id", "order_id", "item_id"}
}

// test items to insert and then select
var (
	items []*tItem = []*tItem{
		{ID: 1, Name: "foo"},
		{ID: 2, Name: "bar"},
		{ID: 3, Name: "test"},
	}
	orders []*tOrder = []*tOrder{
		{ID: 1, Date: "2024-01-01"},
		{ID: 2, Date: "2024-01-02"},
		{ID: 3, Date: "2024-01-03"},
		{ID: 4, Date: "2024-01-04"},
	}
	order_items []*tOrderItem = []*tOrderItem{
		{ID: 1, OrderID: 1, ItemID: 1},
		{ID: 2, OrderID: 2, ItemID: 2},
		{ID: 3, OrderID: 2, ItemID: 3},
		{ID: 4, OrderID: 3, ItemID: 3},
		{ID: 5, OrderID: 4, ItemID: 3},
	}
)
var selectAll string = `
SELECT
	orders.*,
	json_group_array(json_object('id', items.id,'name', items.name)) as items
FROM orders
LEFT JOIN order_items on order_items.order_id = orders.id
LEFT JOIN items on items.id = order_items.item_id
GROUP BY orders.id
ORDER BY orders.date ASC;
`

func setupJoinTestDB(ctx context.Context, path string) (adaptor *adaptors.Sqlite, err error) {

	if adaptor, err = adaptors.NewSqlite(path, false); err != nil {
		return
	}
	// create and insert items
	if _, err = crud.CreateTable(ctx, adaptor, &tItem{}); err != nil {
		return
	}
	if _, err = crud.Insert(ctx, adaptor, &tItem{}, items...); err != nil {
		return
	}
	// create and insert orders
	if _, err = crud.CreateTable(ctx, adaptor, &tOrder{}); err != nil {
		return
	}
	if _, err = crud.Insert(ctx, adaptor, &tOrder{}, orders...); err != nil {
		return
	}
	// create and insert order_items
	if _, err = crud.CreateTable(ctx, adaptor, &tOrderItem{}); err != nil {
		return
	}
	if _, err = crud.Insert(ctx, adaptor, &tOrderItem{}, order_items...); err != nil {
		return
	}
	return
}

func TestAdaptorSelectWithJoins(t *testing.T) {
	fakerextras.AddProviders()
	var (
		err     error
		adaptor *adaptors.Sqlite
		ctx     context.Context = context.Background()
		dir     string          = t.TempDir()
		path    string          = filepath.Join(dir, "test.db")
	)

	adaptor, err = setupJoinTestDB(ctx, path)
	defer adaptor.DB().Close()
	if err != nil {
		t.Fatalf(err.Error())
	}
	results, err := crud.Select[*tOrder](ctx, adaptor, selectAll, nil)
	// check length
	if len(results) != len(orders) {
		t.Errorf("incorrect number of orders found - expected [%d] actual [%v]", len(orders), len(results))
	}
	// check content
	for _, order := range results {
		var expected = 0
		var actual = len(*order.Items)
		for _, o := range order_items {
			if o.OrderID == order.ID {
				expected += 1
			}
		}
		if expected != actual {
			t.Errorf("incorrect number of items for order")
		}
	}

}
