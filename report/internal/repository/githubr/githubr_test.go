package githubr

import (
	"opg-reports/report/internal/utils"

	"github.com/google/go-github/v75/github"
)

// check interfaces are correct
var _ RepositoryReleases = &Repository{}

// simple test release list
var simpleRelID = utils.Ptr[int64](900000)
var simpleReleases = []*github.RepositoryRelease{
	{
		ID:         utils.Ptr[int64](90001),
		Name:       utils.Ptr("previous prerelease"),
		Draft:      utils.Ptr(false),
		Prerelease: utils.Ptr(true),
		TagName:    utils.Ptr("v1.0.1-prerelease"),
		Assets: []*github.ReleaseAsset{
			{
				ID:          utils.Ptr[int64](2001),
				Name:        utils.Ptr("asset.txt"),
				ContentType: utils.Ptr("text/plain"),
			},
			{
				ID:          utils.Ptr[int64](2002),
				Name:        utils.Ptr("asset.json"),
				ContentType: utils.Ptr("text/jsin"),
			},
		},
	},
	{
		ID:         simpleRelID,
		Name:       utils.Ptr("active release with assets"),
		Draft:      utils.Ptr(false),
		Prerelease: utils.Ptr(false),
		TagName:    utils.Ptr("v1.0.1"),
		Assets: []*github.ReleaseAsset{
			{
				ID:          utils.Ptr[int64](1001),
				Name:        utils.Ptr("asset.txt"),
				ContentType: utils.Ptr("text/plain"),
			},
			{
				ID:          utils.Ptr[int64](1002),
				Name:        utils.Ptr("asset.json"),
				ContentType: utils.Ptr("text/plain"),
			},
		},
	},
	{
		ID:         utils.Ptr[int64](89000),
		Name:       utils.Ptr("previous draft"),
		Draft:      utils.Ptr(true),
		Prerelease: utils.Ptr(false),
		TagName:    utils.Ptr("v1.0.1-draft"),
		Assets: []*github.ReleaseAsset{
			{
				ID:          utils.Ptr[int64](3001),
				Name:        utils.Ptr("asset.txt"),
				ContentType: utils.Ptr("text/plain"),
			},
			{
				ID:          utils.Ptr[int64](3002),
				Name:        utils.Ptr("asset.json"),
				ContentType: utils.Ptr("text/plain"),
			},
		},
	},
	{
		ID:         utils.Ptr[int64](80001),
		Name:       utils.Ptr("previous release"),
		Draft:      utils.Ptr(false),
		Prerelease: utils.Ptr(false),
		TagName:    utils.Ptr("v1.0.0"),
	},
}
