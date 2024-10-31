package lib

import (
	"github.com/ministryofjustice/opg-reports/pkg/govukassets"
)

// DownloadGovUKAssets fetches the assets from govuk front end
// and moves it to directory
func DownloadGovUKAssets(directory string) (err error) {
	var frontEnd = govukassets.FrontEnd()
	defer frontEnd.Close()

	_, err = frontEnd.Do(directory)
	return

}
