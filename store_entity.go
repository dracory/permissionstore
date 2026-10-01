package permissionstore

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	contractsorm "github.com/dracory/neat/contracts/database/orm"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

type entityPermissionRow struct {
	ID            string     `db:"id"`
	EntityType    string     `db:"entity_type"`
	EntityID      string     `db:"entity_id"`
	PermissionID  string     `db:"permission_id"`
	Metas         string     `db:"metas"`
	Memo          string     `db:"memo"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	SoftDeletedAt *time.Time `db:"soft_deleted_at"`
}

func (store *store) EntityPermissionCount(ctx context.Context, options EntityPermissionQueryInterface) (int64, error) {
	if options == nil {
		return -1, errors.New("at entityPermission count > entityPermission query is nil")
	}

	options.SetCountOnly(true)

	q, err := store.entityPermissionSelectQuery(options)
	if err != nil {
		return -1, err
	}

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		sqlStr := q.ToRawSql().Get(nil)
		mapped, err := database.SelectToMapString(qCtx, "SELECT COUNT(*) as count FROM ("+sqlStr+") as t")
		if err != nil || len(mapped) < 1 {
			return -1, err
		}
		countStr := mapped[0]["count"]
		return strconv.ParseInt(countStr, 10, 64)
	}

	count, err := q.CountAsVar()
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (store *store) EntityPermissionCreate(ctx context.Context, entityPermission EntityPermissionInterface) error {
	if entityPermission == nil {
		return errors.New("permissionstore > EntityPermissionCreate. entityPermission is nil")
	}

	if entityPermission.PermissionID() == "" {
		return errors.New("permissionstore > EntityPermissionCreate. entityPermission permissionID is empty")
	}

	if entityPermission.EntityID() == "" {
		return errors.New("permissionstore > EntityPermissionCreate. entityPermission entityID is empty")
	}

	if entityPermission.EntityType() == "" {
		return errors.New("permissionstore > EntityPermissionCreate. entityPermission entityType is empty")
	}

	entityPermissionExists, err := store.EntityPermissionFindByEntityAndPermission(
		ctx,
		entityPermission.EntityType(),
		entityPermission.EntityID(),
		entityPermission.PermissionID(),
	)

	if err != nil {
		return err
	}

	if entityPermissionExists != nil {
		return errors.New("permissionstore > EntityPermissionCreate. entityPermission with the same entityType-entityID-permissionID combination already exists")
	}

	entityPermission.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	entityPermission.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := entityPermission.Data()

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		mapString := make(map[string]string)
		for k, v := range data {
			mapString[k] = utils.ToString(v)
		}
		sqlStr := sb.NewBuilder(store.dbDriverName).Table(store.entityPermissionTableName).Insert(mapString)
		_, err := database.Execute(qCtx, sqlStr)
		if err != nil {
			return err
		}
		entityPermission.MarkAsNotDirty()
		return nil
	}

	insertData := make(map[string]interface{})
	for k, v := range data {
		insertData[k] = v
	}

	q := store.db.Query().Table(store.entityPermissionTableName)

	err = q.Create(insertData)
	if err != nil {
		return err
	}

	entityPermission.MarkAsNotDirty()

	return nil
}

func (store *store) EntityPermissionDelete(ctx context.Context, entityPermission EntityPermissionInterface) error {
	if entityPermission == nil {
		return errors.New("entityPermission is nil")
	}

	return store.EntityPermissionDeleteByID(ctx, entityPermission.ID())
}

func (store *store) EntityPermissionDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("entityPermission id is empty")
	}

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		sqlStr := sb.NewBuilder(store.dbDriverName).Table(store.entityPermissionTableName).Where(sb.Where{Column: COLUMN_ID, Operator: "=", Value: id}).Delete()
		_, err := database.Execute(qCtx, sqlStr)
		return err
	}

	q := store.db.Query().Table(store.entityPermissionTableName).Where(COLUMN_ID+" = ?", id)

	_, err := q.Delete()
	return err
}

func (store *store) EntityPermissionFindByEntityAndPermission(
	ctx context.Context,
	entityType string,
	entityID string,
	permissionID string,
) (entityPermission EntityPermissionInterface, err error) {
	if entityType == "" {
		return nil, errors.New("EntityPermissionFindByEntityAndPermission entityType is empty")
	}

	if entityID == "" {
		return nil, errors.New("EntityPermissionFindByEntityAndPermission entityID is empty")
	}

	if permissionID == "" {
		return nil, errors.New("EntityPermissionFindByEntityAndPermission permissionID is empty")
	}

	query := NewEntityPermissionQuery().
		SetEntityType(entityType).
		SetEntityID(entityID).
		SetPermissionID(permissionID).
		SetLimit(1)

	list, err := store.EntityPermissionList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) EntityPermissionFindByID(ctx context.Context, id string) (entityPermission EntityPermissionInterface, err error) {
	if id == "" {
		return nil, errors.New("entityPermission id is empty")
	}

	query := NewEntityPermissionQuery().SetID(id).SetLimit(1)

	list, err := store.EntityPermissionList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) EntityPermissionList(ctx context.Context, query EntityPermissionQueryInterface) ([]EntityPermissionInterface, error) {
	if query == nil {
		return []EntityPermissionInterface{}, errors.New("at entityPermission list > entityPermission query is nil")
	}

	q, err := store.entityPermissionSelectQuery(query)
	if err != nil {
		return []EntityPermissionInterface{}, err
	}

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		sqlStr := q.ToRawSql().Get(nil)
		modelMaps, err := database.SelectToMapString(qCtx, sqlStr)
		if err != nil {
			return []EntityPermissionInterface{}, err
		}
		list := []EntityPermissionInterface{}
		for _, modelMap := range modelMaps {
			model := NewEntityPermissionFromExistingData(modelMap)
			list = append(list, model)
		}
		return list, nil
	}

	var rows []entityPermissionRow
	err = q.Get(&rows)
	if err != nil {
		return []EntityPermissionInterface{}, err
	}

	list := make([]EntityPermissionInterface, 0, len(rows))
	for _, r := range rows {
		ep := NewEntityPermission()
		ep.SetID(r.ID)
		ep.SetEntityType(r.EntityType)
		ep.SetEntityID(r.EntityID)
		ep.SetPermissionID(r.PermissionID)
		if r.Metas != "" {
			var metas map[string]string
			if json.Unmarshal([]byte(r.Metas), &metas) == nil {
				_ = ep.SetMetas(metas)
			}
		}
		ep.SetMemo(r.Memo)
		if r.CreatedAt != nil {
			ep.SetCreatedAt(carbon.CreateFromStdTime(*r.CreatedAt).ToDateTimeString())
		}
		if r.UpdatedAt != nil {
			ep.SetUpdatedAt(carbon.CreateFromStdTime(*r.UpdatedAt).ToDateTimeString())
		}
		if r.SoftDeletedAt != nil {
			ep.SetSoftDeletedAt(carbon.CreateFromStdTime(*r.SoftDeletedAt).ToDateTimeString())
		}
		ep.MarkAsNotDirty()
		list = append(list, ep)
	}

	return list, nil
}

func (store *store) EntityPermissionSoftDelete(ctx context.Context, entityPermission EntityPermissionInterface) error {
	if entityPermission == nil {
		return errors.New("at entityPermission soft delete > entityPermission is nil")
	}

	entityPermission.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.EntityPermissionUpdate(ctx, entityPermission)
}

func (store *store) EntityPermissionSoftDeleteByID(ctx context.Context, id string) error {
	entityPermission, err := store.EntityPermissionFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.EntityPermissionSoftDelete(ctx, entityPermission)
}

func (store *store) EntityPermissionUpdate(ctx context.Context, entityPermission EntityPermissionInterface) error {
	if entityPermission == nil {
		return errors.New("at entityPermission update > entityPermission is nil")
	}

	entityPermission.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := entityPermission.DataChanged()

	delete(dataChanged, COLUMN_ID) // ID is not updateable

	if len(dataChanged) < 1 {
		return nil
	}

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		mapString := make(map[string]string)
		for k, v := range dataChanged {
			mapString[k] = utils.ToString(v)
		}
		sqlStr := sb.NewBuilder(store.dbDriverName).Table(store.entityPermissionTableName).Where(sb.Where{Column: COLUMN_ID, Operator: "=", Value: entityPermission.ID()}).Update(mapString)
		_, err := database.Execute(qCtx, sqlStr)
		if err != nil {
			return err
		}
		entityPermission.MarkAsNotDirty()
		return nil
	}

	updateData := make(map[string]interface{})
	for k, v := range dataChanged {
		updateData[k] = v
	}

	q := store.db.Query().Table(store.entityPermissionTableName).Where(COLUMN_ID+" = ?", entityPermission.ID())

	_, err := q.Update(updateData)
	if err != nil {
		return err
	}

	entityPermission.MarkAsNotDirty()

	return nil
}

func (store *store) entityPermissionSelectQuery(options EntityPermissionQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		return nil, errors.New("entityPermission options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.entityPermissionTableName)

	if len(options.Columns()) > 0 {
		cols := make([]interface{}, len(options.Columns()))
		for i, c := range options.Columns() {
			cols[i] = c
		}
		q = q.Select(cols[0], cols[1:]...)
	}

	if options.HasEntityID() {
		q = q.Where(COLUMN_ENTITY_ID+" = ?", options.EntityID())
	}

	if options.HasEntityType() {
		q = q.Where(COLUMN_ENTITY_TYPE+" = ?", options.EntityType())
	}

	if options.HasID() {
		q = q.Where(COLUMN_ID+" = ?", options.ID())
	}

	if options.HasIDIn() {
		ids := options.IDIn()
		if len(ids) > 0 {
			inClause := COLUMN_ID + " IN ("
			placeholders := make([]interface{}, 0, len(ids))
			for i, id := range ids {
				if i > 0 {
					inClause += ", "
				}
				inClause += "?"
				placeholders = append(placeholders, id)
			}
			inClause += ")"
			q = q.Where(inClause, placeholders...)
		} else {
			q = q.Where("1 = 0")
		}
	}

	if options.HasPermissionID() {
		q = q.Where(COLUMN_PERMISSION_ID+" = ?", options.PermissionID())
	}

	if options.HasCreatedAtGte() {
		q = q.Where(COLUMN_CREATED_AT+" >= ?", options.CreatedAtGte())
	}

	if options.HasCreatedAtLte() {
		q = q.Where(COLUMN_CREATED_AT+" <= ?", options.CreatedAtLte())
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(int(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(int(options.Offset()))
		}
	}

	if options.HasOrderBy() {
		sortDir := lo.Ternary(options.HasSortDirection(), options.SortDirection(), sb.DESC)
		q = q.OrderBy(options.OrderBy(), sortDir)
	}

	if !options.SoftDeletedIncluded() {
		q = q.Where("("+COLUMN_SOFT_DELETED_AT+" > ? OR "+COLUMN_SOFT_DELETED_AT+" IS NULL OR "+COLUMN_SOFT_DELETED_AT+" = '')", carbon.Now(carbon.UTC).ToDateTimeString())
	}

	return q, nil
}
