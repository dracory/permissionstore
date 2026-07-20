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

type permission struct {
	dataobject.DataObject
}

var _ PermissionInterface = (*permission)(nil)

// == CONSTRUCTORS ============================================================

func NewPermission() PermissionInterface {
	o := (&permission{}).
		SetID(uid.HumanUid()).
		SetStatus(PERMISSION_STATUS_INACTIVE).
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

func NewPermissionFromExistingData(data map[string]string) PermissionInterface {
	o := &permission{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *permission) IsActive() bool {
	return o.Status() == PERMISSION_STATUS_ACTIVE
}

func (o *permission) IsSoftDeleted() bool {
	sd := o.SoftDeletedAtCarbon()
	now := carbon.Now(carbon.UTC)
	return sd.Compare("<", now)
}

func (o *permission) IsInactive() bool {
	return o.Status() == PERMISSION_STATUS_INACTIVE
}

// == SETTERS AND GETTERS =====================================================

func (o *permission) CreatedAt() string {
	return o.Get(COLUMN_CREATED_AT)
}

func (o *permission) CreatedAtCarbon() carbon.Carbon {
	return *carbon.Parse(o.CreatedAt(), carbon.UTC)
}

func (o *permission) SetCreatedAt(createdAt string) PermissionInterface {
	o.Set(COLUMN_CREATED_AT, createdAt)
	return o
}

func (o *permission) Handle() string {
	return o.Get(COLUMN_HANDLE)
}

func (o *permission) SetHandle(handle string) PermissionInterface {
	o.Set(COLUMN_HANDLE, handle)
	return o
}

func (o *permission) ID() string {
	return o.Get(COLUMN_ID)
}

func (o *permission) SetID(id string) PermissionInterface {
	o.Set(COLUMN_ID, id)
	return o
}

func (o *permission) Memo() string {
	return o.Get(COLUMN_MEMO)
}

func (o *permission) SetMemo(memo string) PermissionInterface {
	o.Set(COLUMN_MEMO, memo)
	return o
}

func (o *permission) Metas() (map[string]string, error) {
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

func (o *permission) Meta(name string) string {
	metas, err := o.Metas()

	if err != nil {
		return ""
	}

	if value, exists := metas[name]; exists {
		return value
	}

	return ""
}

func (o *permission) SetMeta(name, value string) error {
	return o.UpsertMetas(map[string]string{name: value})
}

// SetMetas stores metas as json string
// Warning: it overwrites any existing metas
func (o *permission) SetMetas(metas map[string]string) error {
	mapString, err := utils.ToJSON(metas)
	if err != nil {
		return err
	}
	o.Set(COLUMN_METAS, mapString)
	return nil
}

func (o *permission) UpsertMetas(metas map[string]string) error {
	currentMetas, err := o.Metas()

	if err != nil {
		return err
	}

	for k, v := range metas {
		currentMetas[k] = v
	}

	return o.SetMetas(currentMetas)
}

func (o *permission) SoftDeletedAt() string {
	return o.Get(COLUMN_SOFT_DELETED_AT)
}

func (o *permission) SoftDeletedAtCarbon() carbon.Carbon {
	return *carbon.Parse(o.SoftDeletedAt(), carbon.UTC)
}

func (o *permission) SetSoftDeletedAt(deletedAt string) PermissionInterface {
	o.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return o
}

func (o *permission) Status() string {
	return o.Get(COLUMN_STATUS)
}

func (o *permission) SetStatus(status string) PermissionInterface {
	o.Set(COLUMN_STATUS, status)
	return o
}

func (o *permission) Title() string {
	return o.Get(COLUMN_TITLE)
}

func (o *permission) SetTitle(title string) PermissionInterface {
	o.Set(COLUMN_TITLE, title)
	return o
}

func (o *permission) UpdatedAt() string {
	return o.Get(COLUMN_UPDATED_AT)
}

func (o *permission) UpdatedAtCarbon() carbon.Carbon {
	return *carbon.Parse(o.Get(COLUMN_UPDATED_AT), carbon.UTC)
}

func (o *permission) SetUpdatedAt(updatedAt string) PermissionInterface {
	o.Set(COLUMN_UPDATED_AT, updatedAt)
	return o
}
