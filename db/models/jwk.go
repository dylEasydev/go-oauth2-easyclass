package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/go-jose/go-jose/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JWKey jose.JSONWebKey

func (j *JWKey) Scan(value any) error {
	if value == nil {
		*j = JWKey(jose.JSONWebKey{})
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("erreur de scanning jwk value :%v", value)
	}

	return j.UnmarshalJSON(bytes)
}

func (j JWKey) Value() (driver.Value, error) {
	return j.MarshalJSON()
}

func (j JWKey) MarshalJSON() ([]byte, error) {
	return jose.JSONWebKey(j).MarshalJSON()
}

// UnmarshalJSON delegates to go-jose's JSON unmarshalling for JSONWebKey.
func (j *JWKey) UnmarshalJSON(b []byte) error {
	var k jose.JSONWebKey
	if err := k.UnmarshalJSON(b); err != nil {
		return err
	}
	*j = JWKey(k)
	return nil
}

func (JWKey) GormDataType() string {
	return "JSON"
}

func (JWKey) GormDBDataTypes(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "JSONB"
	}

	return ""
}
