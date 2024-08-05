package tablequery

import (
	"context"

	"github.com/doublecloud/tross/transfer_manager/go/pkg/abstract"
)

// StorageTableQueryable is storage with table query loading
type StorageTableQueryable interface {
	abstract.Storage

	LoadQueryTable(ctx context.Context, table TableQuery, pusher abstract.Pusher) error
}
