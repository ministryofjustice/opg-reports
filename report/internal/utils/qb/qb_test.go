package qb

import (
	"opg-reports/report/internal/utils/debugger"
	"testing"
)

type variant struct {
	Request  map[string]string
	Expected map[Type][]string
}

type tester struct {
	QB       *Builder
	Variants []*variant
}

var tests = []*tester{

	{
		QB: &Builder{
			From:  "infracosts as base",
			Joins: []string{"LEFT JOIN accounts ON accounts.id = base.account_id"},
			Segments: map[string][]*Segment{
				"_default": {
					{Type: SELECT, Stmt: `strftime("%Y-%m", base.date) as date`},
					{Type: SELECT, Stmt: `CAST(COALESCE(SUM(cost), 0) as REAL) as cost`},
					{Type: WHERE, Stmt: `base.service != 'Tax'`},
					{Type: WHERE, Stmt: `strftime("%Y-%m", base.date) IN (:month)`},
					{Type: GROUPBY, Stmt: `strftime("%Y-%m", base.date)`},
					{Type: ORDERBY, Stmt: `strftime("%Y-%m", base.date) ASC`},
				},
				"team": {
					{Type: SELECT, Stmt: `accounts.team_name as team`},
					{Type: WHERE, Stmt: `accounts.team_name = :team`},
					{Type: GROUPBY, Stmt: `accounts.team_name`},
					{Type: ORDERBY, Stmt: `accounts.team_name ASC`},
				},
				"account": {
					{Type: SELECT, Stmt: `accounts.name as account_name`},
					{Type: SELECT, Stmt: `accounts.id as account_id`},
					{Type: GROUPBY, Stmt: `accounts.name`},
					{Type: WHERE, Stmt: `accounts.name = :account`},
					{Type: ORDERBY, Stmt: `accounts.name ASC`},
				},
				"environment": {
					{Type: SELECT, Stmt: `accounts.environment as account_environment`},
					{Type: WHERE, Stmt: `accounts.environment = :environment`},
					{Type: GROUPBY, Stmt: `accounts.environment`},
					{Type: ORDERBY, Stmt: `accounts.environment ASC`},
				},
				"service": {
					{Type: SELECT, Stmt: `base.service as service`},
					{Type: WHERE, Stmt: `base.service = :service`},
					{Type: GROUPBY, Stmt: `base.service`},
					{Type: ORDERBY, Stmt: `base.service ASC`},
					{Type: HAVING, Stmt: `base.cost > 0`},
				},
			},
		},
		Variants: []*variant{
			// simple filter on the just month
			{
				Request: map[string]string{
					"date_range": "2025-11..2026-02",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
					},
				},
			},
			// filter on the month and group by the team name
			{
				Request: map[string]string{
					"date_range": "2025-11..2026-02",
					"team":       "true",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
						`accounts.team_name as team`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
						`accounts.team_name`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
						`accounts.team_name ASC`,
					},
				},
			},
			// filter on the month and the team
			{
				Request: map[string]string{
					"date_range": "2025-11..2026-02",
					"team":       "sirius",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
						`accounts.team_name as team`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
						`accounts.team_name = :team`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
						`accounts.team_name`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
						`accounts.team_name ASC`,
					},
				},
			},
			// filter on the month and the team and group by the account
			{
				Request: map[string]string{
					"date_range": "2025-11..2026-02",
					"team":       "sirius",
					"account":    "true",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
						`accounts.team_name as team`,
						`accounts.name as account_name`,
						`accounts.id as account_id`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
						`accounts.team_name = :team`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
						`accounts.team_name`,
						`accounts.name`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
						`accounts.team_name ASC`,
						`accounts.name ASC`,
					},
				},
			},
			// filter on the month, team and account
			{
				Request: map[string]string{
					"date_range": "2025-11..2026-02",
					"team":       "sirius",
					"account":    "sirius production",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
						`accounts.team_name as team`,
						`accounts.name as account_name`,
						`accounts.id as account_id`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
						`accounts.team_name = :team`,
						`accounts.name = :account`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
						`accounts.team_name`,
						`accounts.name`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
						`accounts.team_name ASC`,
						`accounts.name ASC`,
					},
				},
			},
			// add in environment
			{
				Request: map[string]string{
					"date_range":  "2025-11..2026-02",
					"team":        "sirius",
					"account":     "sirius production",
					"environment": "true",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
						`accounts.team_name as team`,
						`accounts.name as account_name`,
						`accounts.id as account_id`,
						`accounts.environment as account_environment`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
						`accounts.team_name = :team`,
						`accounts.name = :account`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
						`accounts.team_name`,
						`accounts.name`,
						`accounts.environment`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
						`accounts.team_name ASC`,
						`accounts.name ASC`,
						`accounts.environment ASC`,
					},
				},
			},
			// filter by environment
			{
				Request: map[string]string{
					"date_range":  "2025-11..2026-02",
					"team":        "sirius",
					"account":     "sirius production",
					"environment": "production",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
						`accounts.team_name as team`,
						`accounts.name as account_name`,
						`accounts.id as account_id`,
						`accounts.environment as account_environment`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
						`accounts.team_name = :team`,
						`accounts.name = :account`,
						`accounts.environment = :environment`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
						`accounts.team_name`,
						`accounts.name`,
						`accounts.environment`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
						`accounts.team_name ASC`,
						`accounts.name ASC`,
						`accounts.environment ASC`,
					},
				},
			},
			// add service
			{
				Request: map[string]string{
					"date_range":  "2025-11..2026-02",
					"team":        "sirius",
					"account":     "sirius production",
					"environment": "production",
					"service":     "true",
				},
				Expected: map[Type][]string{
					SELECT: []string{
						`strftime("%Y-%m", base.date) as date`,
						`CAST(COALESCE(SUM(cost), 0) as REAL) as cost`,
						`accounts.team_name as team`,
						`accounts.name as account_name`,
						`accounts.id as account_id`,
						`accounts.environment as account_environment`,
						`base.service as service`,
					},
					WHERE: []string{
						`base.service != 'Tax'`,
						`strftime("%Y-%m", base.date) IN (:month)`,
						`accounts.team_name = :team`,
						`accounts.name = :account`,
						`accounts.environment = :environment`,
					},
					GROUPBY: []string{
						`strftime("%Y-%m", base.date)`,
						`accounts.team_name`,
						`accounts.name`,
						`accounts.environment`,
						`base.service`,
					},
					ORDERBY: []string{
						`strftime("%Y-%m", base.date) ASC`,
						`accounts.team_name ASC`,
						`accounts.name ASC`,
						`accounts.environment ASC`,
						`base.service ASC`,
					},
					HAVING: []string{
						`base.cost > 0`,
					},
				},
			},
		},
	},
}

func TestUtilsBuilder(t *testing.T) {

	for _, test := range tests {
		var q = test.QB
		for i, variant := range test.Variants {
			_, blocks := q.FromRequest(variant.Request)
			if len(blocks) != len(variant.Expected) {
				t.Errorf("different number of sql segments returned. expected [%d] actual [%v].", len(variant.Expected), len(blocks))
			}
			// loop over expected chunks
			for ty, expected := range variant.Expected {
				// make sure they match exactly the expected values
				if len(blocks[ty]) != len(expected) {
					t.Errorf("[%d] different number of segments in [%s] for statement. expected [%d] actual [%v]\n[%s].",
						i, ty, len(expected), len(blocks[ty]), debugger.DumpStr(blocks[ty]))

				}
				// look for every expected block and find it in the generated set
				for _, ex := range expected {
					var found = false
					for _, block := range blocks[ty] {
						if block == ex {
							found = true
						}
					}
					if !found {
						t.Errorf("expected segment [%s] block [%s] not found in:\n[%s]", ty, ex, debugger.DumpStr(blocks[ty]))
					}
				}
			}
		}

	}

}
