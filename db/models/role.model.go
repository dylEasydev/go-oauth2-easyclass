package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/utils"
	"github.com/dylEasydev/go-oauth2-easyclass/validators"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScopeJSON struct {
	ScopeName     string `json:"scopeName"`
	ScopeDescript string `json:"scopeDescript"`
}
type ScopeData struct {
	Data []ScopeJSON `json:"data"`
}

// structure des role d'utilisateur (admin , teacher , student ...)
type Role struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	RoleName     string    `gorm:"column:rolename;not null;unique" validate:"required,rowallowed"`
	RoleDescript string    `gorm:"column:roledescript"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Scopes []*Scope `gorm:"many2many:authpermission;"`
}

// implementation de l'interface Tabler
func (Role) TableName() string {
	return "roles"
}

// validation de l'entrée
func (role *Role) BeforeSave(tx *gorm.DB) error {
	return validators.ValidateStruct(role)
}

// fonction ajout des permission d'un role
func (role *Role) AddScope(tx *gorm.DB) (err error) {
	// construction du chemin directeur .
	name := fmt.Sprintf("ressources/scope_%s", strings.ToLower(role.RoleName))

	//lecture des permissions
	data, err := utils.ReadJSON[ScopeData](name)
	if err != nil {
		return fmt.Errorf("lecture scopes pour role %s : %w", role.RoleName, err)
	}

	// regroupement des nom de scope dans une tranche
	//pour faciliter la reherche
	var sliceName []string = make([]string, len(data.Data))
	for i, elem := range data.Data {
		sliceName[i] = elem.ScopeName
	}

	//recherche des permission correspondant en BD
	ctx := context.Background()

	scopes, err := gorm.G[Scope](tx).Where("scopename IN ?", sliceName).Find(ctx)
	if err != nil {
		return
	}

	//remplacement des associations pou se role
	if len(scopes) > 0 {
		if err = tx.Model(role).Association("Scopes").Replace(&scopes); err != nil {
			return fmt.Errorf("erreur remplacement scopes : %w", err)
		}
	}

	return nil
}

// hooks après la sauvegarde des roles
// ajouter les scopes
func (role *Role) AfterSave(tx *gorm.DB) error {
	if err := role.AddScope(tx); err != nil {
		fmt.Printf("Erreur lors de la mise à jour des scopes: %v\n", err)
	}
	return nil
}
