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

type permissionRow struct {
	ID            string     `db:"id"`
	Status        string     `db:"status"`
	Handle        string     `db:"handle"`
	Title         string     `db:"title"`
	Metas         string     `db:"metas"`
	Memo          string     `db:"memo"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	SoftDeletedAt *time.Time `db:"soft_deleted_at"`
}

func (store *store) PermissionCount(ctx context.Context, options PermissionQueryInterface) (int64, error) {
	if options == nil {
		return -1, errors.New("at permission count > permission query is nil")
	}

	options.SetCountOnly(true)

	q, err := store.permissionSelectQuery(options)
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

func (store *store) PermissionCreate(ctx context.Context, permission PermissionInterface) error {
	if permission == nil {
		return errors.New("permission is nil")
	}

	permission.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))
	permission.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	data := permission.Data()

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		mapString := make(map[string]string)
		for k, v := range data {
			mapString[k] = utils.ToString(v)
		}
		sqlStr := sb.NewBuilder(store.dbDriverName).Table(store.permissionTableName).Insert(mapString)
		_, err := database.Execute(qCtx, sqlStr)
		if err != nil {
			return err
		}
		permission.MarkAsNotDirty()
		return nil
	}

	insertData := make(map[string]interface{})
	for k, v := range data {
		insertData[k] = v
	}

	q := store.db.Query().Table(store.permissionTableName)

	err := q.Create(insertData)
	if err != nil {
		return err
	}

	permission.MarkAsNotDirty()

	return nil
}

func (store *store) PermissionDelete(ctx context.Context, permission PermissionInterface) error {
	if permission == nil {
		return errors.New("permission is nil")
	}

	return store.PermissionDeleteByID(ctx, permission.ID())
}

func (store *store) PermissionDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("permission id is empty")
	}

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		sqlStr := sb.NewBuilder(store.dbDriverName).Table(store.permissionTableName).Where(sb.Where{Column: COLUMN_ID, Operator: "=", Value: id}).Delete()
		_, err := database.Execute(qCtx, sqlStr)
		return err
	}

	q := store.db.Query().Table(store.permissionTableName).Where(COLUMN_ID+" = ?", id)

	_, err := q.Delete()
	return err
}

func (store *store) PermissionFindByHandle(ctx context.Context, handle string) (permission PermissionInterface, err error) {
	if handle == "" {
		return nil, errors.New("permission handle is empty")
	}

	query := NewPermissionQuery().SetHandle(handle).SetLimit(1)

	list, err := store.PermissionList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) PermissionFindByID(ctx context.Context, id string) (permission PermissionInterface, err error) {
	if id == "" {
		return nil, errors.New("permission id is empty")
	}

	query := NewPermissionQuery().SetID(id).SetLimit(1)

	list, err := store.PermissionList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) PermissionList(ctx context.Context, query PermissionQueryInterface) ([]PermissionInterface, error) {
	if query == nil {
		return []PermissionInterface{}, errors.New("at permission list > permission query is nil")
	}

	q, err := store.permissionSelectQuery(query)
	if err != nil {
		return []PermissionInterface{}, err
	}

	qCtx := store.toQuerableContext(ctx)
	if qCtx.IsTx() {
		sqlStr := q.ToRawSql().Get(nil)
		modelMaps, err := database.SelectToMapString(qCtx, sqlStr)
		if err != nil {
			return []PermissionInterface{}, err
		}
		list := []PermissionInterface{}
		for _, modelMap := range modelMaps {
			model := NewPermissionFromExistingData(modelMap)
			list = append(list, model)
		}
		return list, nil
	}

	var rows []permissionRow
	err = q.Get(&rows)
	if err != nil {
		return []PermissionInterface{}, err
	}

	list := make([]PermissionInterface, 0, len(rows))
	for _, r := range rows {
		p := NewPermission()
		p.SetID(r.ID)
		p.SetStatus(r.Status)
		p.SetHandle(r.Handle)
		p.SetTitle(r.Title)
		if r.Metas != "" {
			var metas map[string]string
			if json.Unmarshal([]byte(r.Metas), &metas) == nil {
				_ = p.SetMetas(metas)
			}
		}
		p.SetMemo(r.Memo)
		if r.CreatedAt != nil {
			p.SetCreatedAt(carbon.CreateFromStdTime(*r.CreatedAt).ToDateTimeString())
		}
		if r.UpdatedAt != nil {
			p.SetUpdatedAt(carbon.CreateFromStdTime(*r.UpdatedAt).ToDateTimeString())
		}
		if r.SoftDeletedAt != nil {
			p.SetSoftDeletedAt(carbon.CreateFromStdTime(*r.SoftDeletedAt).ToDateTimeString())
		}
		p.MarkAsNotDirty()
		list = append(list, p)
	}

	return list, nil
}

func (store *store) PermissionSoftDelete(ctx context.Context, permission PermissionInterface) error {
	if permission == nil {
		return errors.New("at permission soft delete > permission is nil")
	}

	permission.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.PermissionUpdate(ctx, permission)
}

func (store *store) PermissionSoftDeleteByID(ctx context.Context, id string) error {
	permission, err := store.PermissionFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.PermissionSoftDelete(ctx, permission)
}

func (store *store) PermissionUpdate(ctx context.Context, permission PermissionInterface) error {
	if permission == nil {
		return errors.New("at permission update > permission is nil")
	}

	permission.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())

	dataChanged := permission.DataChanged()

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
		sqlStr := sb.NewBuilder(store.dbDriverName).Table(store.permissionTableName).Where(sb.Where{Column: COLUMN_ID, Operator: "=", Value: permission.ID()}).Update(mapString)
		_, err := database.Execute(qCtx, sqlStr)
		if err != nil {
			return err
		}
		permission.MarkAsNotDirty()
		return nil
	}

	updateData := make(map[string]interface{})
	for k, v := range dataChanged {
		updateData[k] = v
	}

	q := store.db.Query().Table(store.permissionTableName).Where(COLUMN_ID+" = ?", permission.ID())

	_, err := q.Update(updateData)
	if err != nil {
		return err
	}

	permission.MarkAsNotDirty()

	return nil
}

func (store *store) permissionSelectQuery(options PermissionQueryInterface) (contractsorm.Query, error) {
	if options == nil {
		return nil, errors.New("permission options is nil")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	q := store.db.Query().Table(store.permissionTableName)

	if len(options.Columns()) > 0 {
		cols := make([]interface{}, len(options.Columns()))
		for i, c := range options.Columns() {
			cols[i] = c
		}
		q = q.Select(cols[0], cols[1:]...)
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

	if options.HasStatus() {
		q = q.Where(COLUMN_STATUS+" = ?", options.Status())
	}

	if options.HasStatusIn() {
		statuses := options.StatusIn()
		if len(statuses) > 0 {
			inClause := COLUMN_STATUS + " IN ("
			placeholders := make([]interface{}, 0, len(statuses))
			for i, status := range statuses {
				if i > 0 {
					inClause += ", "
				}
				inClause += "?"
				placeholders = append(placeholders, status)
			}
			inClause += ")"
			q = q.Where(inClause, placeholders...)
		} else {
			q = q.Where("1 = 0")
		}
	}

	if options.HasHandle() {
		q = q.Where(COLUMN_HANDLE+" = ?", options.Handle())
	}

	if options.HasTitleLike() {
		q = q.Where(COLUMN_TITLE+" LIKE ?", "%"+options.TitleLike()+"%")
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
