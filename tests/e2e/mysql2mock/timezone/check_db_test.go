package nonutf8charset

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/doublecloud/transfer/internal/logger"
	"github.com/doublecloud/transfer/library/go/test/canon"
	"github.com/doublecloud/transfer/pkg/abstract"
	"github.com/doublecloud/transfer/pkg/abstract/coordinator"
	server "github.com/doublecloud/transfer/pkg/abstract/model"
	mysql_storage "github.com/doublecloud/transfer/pkg/providers/mysql"
	"github.com/doublecloud/transfer/pkg/runtime/local"
	"github.com/doublecloud/transfer/tests/helpers"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

const tableName = "__test1"

var (
	db     = os.Getenv("RECIPE_MYSQL_SOURCE_DATABASE")
	source = helpers.WithMysqlInclude(
		helpers.RecipeMysqlSource(),
		[]string{fmt.Sprintf("%s.%s", db, tableName)},
	)
)

func init() {
	source.WithDefaults()
	source.Timezone = "Europe/Moscow"
}

type mockSinker struct {
	pushCallback func(input []abstract.ChangeItem) error
}

func (s *mockSinker) Push(input []abstract.ChangeItem) error {
	return s.pushCallback(input)
}

func (s *mockSinker) Close() error {
	return nil
}

func dummyProgress(current, progress, total uint64) {}

func makeConnConfig() *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.Addr = fmt.Sprintf("%v:%v", source.Host, source.Port)
	cfg.User = source.User
	cfg.Passwd = string(source.Password)
	cfg.DBName = source.Database
	cfg.Net = "tcp"
	return cfg
}

func TestTimeZone(t *testing.T) {
	defer func() {
		require.NoError(t, helpers.CheckConnections(
			helpers.LabeledPort{Label: "Mysql source", Port: source.Port},
		))
	}()

	storage, err := mysql_storage.NewStorage(source.ToStorageParams())
	require.NoError(t, err)

	rowsValues := []any{}

	table := abstract.TableDescription{Name: tableName, Schema: source.Database}
	err = storage.LoadTable(context.Background(), table, func(input []abstract.ChangeItem) error {
		for _, item := range input {
			if item.Kind != "insert" {
				continue
			}
			rowsValues = append(rowsValues, item.ColumnValues)
		}
		return nil
	})
	require.NoError(t, err)

	var sinker mockSinker
	target := server.MockDestination{SinkerFactory: func() abstract.Sinker {
		return &sinker
	}}
	transfer := server.Transfer{
		ID:  "test",
		Src: source,
		Dst: &target,
	}

	fakeClient := coordinator.NewStatefulFakeClient()
	err = mysql_storage.SyncBinlogPosition(source, transfer.ID, fakeClient)
	require.NoError(t, err)

	wrk := local.NewLocalWorker(fakeClient, &transfer, helpers.EmptyRegistry(), logger.Log)

	sinker.pushCallback = func(input []abstract.ChangeItem) error {
		for _, item := range input {
			if item.Kind != "insert" {
				continue
			}
			rowsValues = append(rowsValues, item.ColumnValues)
		}

		if len(rowsValues) >= 4 {
			_ = wrk.Stop()
		}

		return nil
	}

	errCh := make(chan error)
	go func() {
		errCh <- wrk.Run()
	}()

	conn, err := mysql.NewConnector(makeConnConfig())
	require.NoError(t, err)
	db := sql.OpenDB(conn)

	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	require.NoError(t, err)

	_, err = tx.Query("SET SESSION time_zone = '+00:00';")
	require.NoError(t, err)

	_, err = tx.Query(fmt.Sprintf(`
		INSERT INTO %s (ts) VALUES
			('2020-12-23 10:11:12'),
			('2020-12-23 14:15:16');
	`, tableName))
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	require.NoError(t, <-errCh)

	canon.SaveJSON(t, rowsValues)
}