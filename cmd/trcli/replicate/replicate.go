package replicate

import (
	"github.com/doublecloud/transfer/cmd/trcli/config"
	"github.com/doublecloud/transfer/internal/logger"
	"github.com/doublecloud/transfer/library/go/core/metrics/solomon"
	"github.com/doublecloud/transfer/library/go/core/xerrors"
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/doublecloud/transfer/pkg/abstract/coordinator"
	"github.com/doublecloud/transfer/pkg/abstract/model"
	"github.com/doublecloud/transfer/pkg/dataplane/provideradapter"
	"github.com/doublecloud/transfer/pkg/runtime/local"
	"github.com/spf13/cobra"
)

func ReplicateCommand(cp *coordinator.Coordinator, rt abstract.Runtime) *cobra.Command {
	var transferParams string
	replicationCommand := &cobra.Command{
		Use:   "replicate",
		Short: "Start local replication",
		RunE:  replicate(cp, rt, &transferParams),
	}
	replicationCommand.Flags().StringVar(&transferParams, "transfer", "./transfer.yaml", "path to yaml file with transfer configuration")
	return replicationCommand
}

func replicate(cp *coordinator.Coordinator, rt abstract.Runtime, transferYaml *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		transfer, err := config.TransferFromYaml(transferYaml)
		if err != nil {
			return xerrors.Errorf("unable to load transfer: %w", err)
		}
		transfer.Runtime = rt

		return RunReplication(*cp, transfer)
	}
}

func RunReplication(cp coordinator.Coordinator, transfer *model.Transfer) error {
	if err := provideradapter.ApplyForTransfer(transfer); err != nil {
		return xerrors.Errorf("unable to adapt transfer: %w", err)
	}
	for {
		worker := local.NewLocalWorker(cp, transfer, solomon.NewRegistry(solomon.NewRegistryOpts()), logger.Log)
		err := worker.Run()
		if abstract.IsFatal(err) {
			return err
		}
		if err := worker.Stop(); err != nil {
			logger.Log.Warnf("unable to stop worker: %v", err)
		}
		logger.Log.Warnf("worker failed: %v, restart", err)
	}
}