package team

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/go-github/v62/github"
	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/ministryofjustice/opg-reports/report/internal/gh"
	"github.com/ministryofjustice/opg-reports/report/internal/interfaces"
	"github.com/ministryofjustice/opg-reports/report/internal/sqldb"
	"github.com/ministryofjustice/opg-reports/report/internal/utils"
)

type Service[T interfaces.Model] struct {
	ctx   context.Context
	log   *slog.Logger
	conf  *config.Config
	store *sqldb.Repository[T]
}

// Seed is used to insert test data to the team table, so for now
// we create 3 dummy teams
func (self *Service[T]) Seed() (err error) {
	var (
		now   = time.Now().UTC().Format(time.RFC3339)
		seeds = []*sqldb.BoundStatement{
			{Statement: stmtInsert, Data: &Team{Name: "TeamA", CreatedAt: now}},
			{Statement: stmtInsert, Data: &Team{Name: "TeamB", CreatedAt: now}},
			{Statement: stmtInsert, Data: &Team{Name: "TeamC", CreatedAt: now}},
		}
	)
	err = self.store.Insert(seeds...)

	return
}

// Import will fetch published opg-metadata content and convert all 'billing_unit'
// entries within the opg-metadata account list to become teams.
//
// The latest release is downloaded, extracted and the accounts.json used as the
// dataset
//
// Uses the gh repository to fetch the data from github
func (self *Service[T]) Import(gh *gh.Repository) (err error) {
	var (
		asset          *github.ReleaseAsset
		downloadedFile *os.File
		org            string = self.conf.Github.Organisation
		dataRepo       string = "opg-metadata"
		assetName      string = "metadata.tar.gz"
		downloadTo     string = "./team-service-import/" + assetName
		extractTo      string = "./team-service-import/" + dataRepo
	)
	// get the latest relase and the asset details that match the name
	asset, err = gh.GetLatestReleaseAsset(org, dataRepo, assetName, false)
	if err != nil {
		return
	}
	// download this asset
	downloadedFile, err = gh.DownloadReleaseAsset(org, dataRepo, *asset.ID, downloadTo)
	if err != nil {
		return
	}
	defer downloadedFile.Close()

	// now extract the tar.gz
	err = utils.TarGzExtract(downloadedFile, extractTo)
	return
}

// GetAllTeams returns all teams as a slice from the database
// Calls the database
func (self *Service[T]) GetAllTeams() (teams []*Team, err error) {
	var selectStmt = &sqldb.BoundStatement{Statement: stmtSelectAll}
	teams = []*Team{}

	err = self.store.Select(selectStmt)
	// cast the data back to struct
	if err == nil {
		teams = selectStmt.Returned.([]*Team)
	}

	return
}

func NewService[T interfaces.Model](ctx context.Context, log *slog.Logger, conf *config.Config, store *sqldb.Repository[T]) (srv *Service[T], err error) {
	if log == nil {
		return nil, fmt.Errorf("no logger passed for team service")
	}
	if conf == nil {
		return nil, fmt.Errorf("no config passed for team service")
	}
	if store == nil {
		return nil, fmt.Errorf("no repository passed for team service")
	}

	srv = &Service[T]{
		ctx:   ctx,
		log:   log.With("service", "team"),
		conf:  conf,
		store: store,
	}
	return
}
