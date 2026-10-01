package permissionstore

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/dracory/neat"
	"github.com/gouniverse/base/database"
)

// == TYPE ====================================================================

type store struct {
	// permissionTableName is the name of the permission table
	permissionTableName string

	// entityPermissionTableName is the name of the permission entity relation table
	entityPermissionTableName string

	// db is the underlying neat database
	db *neat.Database

	// dbDriverName is the database driver name/type
	dbDriverName string

	// automigrateEnabled enables or disables automigration
	automigrateEnabled bool

	// debugEnabled enables or disables debug mode
	debugEnabled bool

	// sqlLogger is the sql logger used when debug mode is enabled
	sqlLogger *slog.Logger
}

// == INTERFACE ===============================================================

var _ StoreInterface = (*store)(nil) // verify it extends the interface

// PUBLIC METHODS ============================================================

// DB returns the underlying database connection
func (store *store) DB() *sql.DB {
	if store.db == nil {
		return nil
	}
	db, err := store.db.DB()
	if err != nil {
		return nil
	}
	return db
}

// EnableDebug - enables or disables the debug mode
func (st *store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// logSql logs sql to the sql logger, if debug mode is enabled
func (store *store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}

// toQuerableContext converts the context to a QueryableContext
func (store *store) toQuerableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}

	return database.Context(ctx, store.DB())
}
