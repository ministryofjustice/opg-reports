package main

import "github.com/ministryofjustice/opg-reports/servers/sfront/lib"

const assetsDirectory string = "./"

func main() {
	// download the assets
	lib.DownloadGovUKAssets(assetsDirectory)
}
