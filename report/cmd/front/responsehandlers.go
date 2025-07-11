package main

import (
	"opg-reports/report/internal/utils"
)

// apiResponseTeams captures name only data from the api response
// which is then used for generating the navigation structure
//
// endpoint: `/v1/teams/all`
type apiResponseTeams struct {
	Count int `json:"count,omityempty"`
	Data  []*struct {
		Name string `json:"name"`
	} `json:"data"`
}

// parseAllTeamsForNavigation excludes Legacy & ORG from team listing in the navigation
// for ease
func parseAllTeamsForNavigation(response *apiResponseTeams) (teams []string, err error) {
	teams = []string{}
	for _, team := range response.Data {
		if team.Name != "Legacy" && team.Name != "ORG" {
			teams = append(teams, team.Name)
		}
	}
	return
}

// apiResponseAwsCostsGrouped
type apiResponseAwsCostsGrouped struct {
	Count  int                 `json:"count,omityempty"`
	Dates  []string            `json:"dates,omitempty"`
	Groups []string            `json:"groups,omitempty"`
	Data   []map[string]string `json:"data"`
}

type dataTable struct {
	Body         map[string]map[string]string
	RowHeaders   []string
	DataHeaders  []string
	ExtraHeaders []string
	Header       []string
	Footer       map[string]string
}

// parseAwsCostsGrouped TODO - CLEAN UP!
func parseAwsCostsGrouped(response *apiResponseAwsCostsGrouped) (dt *dataTable, err error) {

	dummyRow := utils.DummyRows(response.Dates, "date")
	allData := append(response.Data, dummyRow...)
	possibles, _ := utils.PossibleCombinationsAsKeys(allData, response.Groups)
	table := utils.SkeletonTable(possibles, response.Dates)
	body := utils.PopulateTable(response.Data, table, response.Groups, "date", "cost")

	utils.AddRowTotals(body, response.Groups, "total")
	utils.AddColumnsToRows(body, "trend")

	extraheaders := []string{"trend", "total"}
	sumCols := append(response.Dates, "total")
	extratotals := append(response.Groups, "trend")

	dt = &dataTable{
		Body:         body,
		RowHeaders:   response.Groups,
		DataHeaders:  response.Dates,
		ExtraHeaders: extraheaders,
		Header:       utils.TableHeaderRow(response.Groups, response.Dates, extraheaders),
		Footer:       utils.ColumnTotals(body, sumCols, extratotals...),
	}

	utils.Debug(dt)

	return
}
