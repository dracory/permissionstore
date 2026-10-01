package permissionstore

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/dracory/neat"
	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	// PermissionTableName is the name of the permission table
	PermissionTableName string

	// EntityPermissionTableName is the name of the entity to permission relation table
	EntityPermissionTableName string

	// DB is the underlying database connection
	DB *sql.DB

	// DbDriverName is the database driver name/type
	DbDriverName string

	// AutomigrateEnabled indicates whether to automatically migrate the database
	AutomigrateEnabled bool

	// DebugEnabled enables or disables the debug mode
	DebugEnabled bool

	// SqlLogger is the sql statement logger when debug mode is enabled, defaults to the default logger
	SqlLogger *slog.Logger
}

// NewStore creates a new block store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.PermissionTableName == "" {
		return nil, errors.New("permission store: PermissionTableName is required")
	}

	if opts.EntityPermissionTableName == "" {
		return nil, errors.New("permission store: EntityPermissionTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("permission store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	if opts.SqlLogger == nil {
		opts.SqlLogger = slog.Default()
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	store := &store{
		permissionTableName:       opts.PermissionTableName,
		entityPermissionTableName: opts.EntityPermissionTableName,
		automigrateEnabled:        opts.AutomigrateEnabled,
		db:                        neatDB,
		dbDriverName:              opts.DbDriverName,
		debugEnabled:              opts.DebugEnabled,
		sqlLogger:                 opts.SqlLogger,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
