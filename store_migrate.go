package permissionstore

import (
	"errors"

	contractsschema "github.com/dracory/neat/contracts/database/schema"
)

// AutoMigrate auto-migrates the database schema using neat
func (store *store) AutoMigrate() error {
	if store.db == nil {
		return errors.New("permissionstore: database is nil")
	}

	if !store.db.Schema().HasTable(store.permissionTableName) {
		err := store.db.Schema().Create(store.permissionTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_STATUS, 40).Default("")
			table.String(COLUMN_HANDLE, 50).Default("")
			table.String(COLUMN_TITLE, 100).Default("")
			table.Text(COLUMN_METAS).Default("")
			table.Text(COLUMN_MEMO).Default("")
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			return err
		}
	}

	if !store.db.Schema().HasTable(store.entityPermissionTableName) {
		err := store.db.Schema().Create(store.entityPermissionTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_ENTITY_TYPE, 80).Default("")
			table.String(COLUMN_ENTITY_ID, 40).Default("")
			table.String(COLUMN_PERMISSION_ID, 40).Default("")
			table.Text(COLUMN_METAS).Default("")
			table.Text(COLUMN_MEMO).Default("")
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
