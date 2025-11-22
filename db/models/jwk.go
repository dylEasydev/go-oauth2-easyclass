package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/go-jose/go-jose/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JWKey jose.JSONWebKey

func (j *JWKey) Scan(value any) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("errur de marshaling value :%v", value)
	}
	result := jose.JSONWebKey{}
	err := json.Unmarshal(bytes, &result)
	*j = JWKey(result)
	return err
}

func (j JWKey) Value() (driver.Value, error) {
	return json.Marshal(j)
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
