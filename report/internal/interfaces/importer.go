package interfaces

import (
	"context"
	"log/slog"

	"github.com/ministryofjustice/opg-reports/report/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ImporterCLICommand func(conf *config.Config, viperConf *viper.Viper) (cmd *cobra.Command)
type ImporterExistingCommand func(ctx context.Context, log *slog.Logger, conf *config.Config) (err error)
