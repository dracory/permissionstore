package permissionstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
	"github.com/gouniverse/utils"
)

// == CLASS ===================================================================

type entityPermission struct {
	dataobject.DataObject
}

var _ EntityPermissionInterface = (*entityPermission)(nil)

// == CONSTRUCTORS ============================================================

func NewEntityPermission() EntityPermissionInterface {
	o := (&entityPermission{}).
		SetID(uid.HumanUid()).
		SetMemo("").
		SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)).
		SetSoftDeletedAt(sb.MAX_DATETIME)

	err := o.SetMetas(map[string]string{})

	if err != nil {
		return o
	}

	return o
}

func NewEntityPermissionFromExistingData(data map[string]string) EntityPermissionInterface {
	o := &entityPermission{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *entityPermission) IsSoftDeleted() bool {
	sd := o.SoftDeletedAtCarbon()
	now := carbon.Now(carbon.UTC)
	return sd.Compare("<", now)
}

// == SETTERS AND GETTERS =====================================================

func (o *entityPermission) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *entityPermission) CreatedAtCarbon() carbon.Carbon {
	return *carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *entityPermission) SetCreatedAt(createdAt string) EntityPermissionInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *entityPermission) EntityType() string {
	return o.Get(COLUMN_ENTITY_TYPE)
}

func (o *entityPermission) SetEntityType(entityType string) EntityPermissionInterface {
	o.Set(COLUMN_ENTITY_TYPE, entityType)
	return o
}

func (o *entityPermission) EntityID() string {
	return o.Get(COLUMN_ENTITY_ID)
}

func (o *entityPermission) SetEntityID(entityID string) EntityPermissionInterface {
	o.Set(COLUMN_ENTITY_ID, entityID)
	return o
}

func (o *entityPermission) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *entityPermission) SetID(id string) EntityPermissionInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *entityPermission) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *entityPermission) SetMemo(memo string) EntityPermissionInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *entityPermission) Metas() (map[string]string, error) {
	metasStr := o.Get(COLUMN_METAS)

	if metasStr == "" {
		metasStr = "{}"
	}

	metasJson, errJson := utils.FromJSON(metasStr, map[string]string{})
	if errJson != nil {
		return map[string]string{}, errJson
	}

	return maputils.MapStringAnyToMapStringString(metasJson.(map[string]any)), nil
}

func (o *entityPermission) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *entityPermission) SetMeta(name, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *entityPermission) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}
	o.Set(COLUMN_METAS, mapString)
	return nil
}

func (o *entityPermission) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *entityPermission) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *entityPermission) SoftDeletedAtCarbon() carbon.Carbon {
	return *carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *entityPermission) SetSoftDeletedAt(deletedAt string) EntityPermissionInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return o
}

func (o *entityPermission) PermissionID() string {
	return o.Get(COLUMN_PERMISSION_ID)
}

func (o *entityPermission) SetPermissionID(roleID string) EntityPermissionInterface {
	o.Set(COLUMN_PERMISSION_ID, roleID)
	return o
}

func (o *entityPermission) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *entityPermission) UpdatedAtCarbon() carbon.Carbon {
	return *carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *entityPermission) SetUpdatedAt(updatedAt string) EntityPermissionInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
