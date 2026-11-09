package database

import (
	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	scopes := []domain.Scope{
		{Name: "events:list"},
		{Name: "events:read"},
		{Name: "seats:list"},
		{Name: "bookings:reserve"},
		{Name: "bookings:confirm"},
		{Name: "bookings:reject"},
	}

	for i := range scopes {
		if err := db.Where(domain.Scope{Name: scopes[i].Name}).FirstOrCreate(&scopes[i]).Error; err != nil {
			return err
		}
	}

	scopeMap := make(map[string]domain.Scope)
	for _, s := range scopes {
		scopeMap[s.Name] = s
	}

	rolesConfig := map[string][]string{
		"user": {
			"events:list",
			"events:read",
			"seats:list",
			"bookings:reserve",
		},
		"admin": {
			"events:list",
			"events:read",
			"seats:list",
			"bookings:reserve",
			"bookings:confirm",
			"bookings:reject",
		},
		"payment_gtw": {
			"bookings:confirm",
			"bookings:reject",
		},
	}

	for roleName, scopeNames := range rolesConfig {
		var role domain.Role
		if err := db.Where(domain.Role{Name: roleName}).FirstOrCreate(&role).Error; err != nil {
			return err
		}

		var roleScopes []domain.Scope
		for _, sName := range scopeNames {
			if s, exists := scopeMap[sName]; exists {
				roleScopes = append(roleScopes, s)
			}
		}

		if err := db.Model(&role).Association("Scopes").Replace(roleScopes); err != nil {
			return err
		}
	}

	roles := []domain.Role{
		{Name: "user"},
		{Name: "admin"},
		{Name: "payment_gtw"},
	}

	for i := range roles {
		if err := db.Where(domain.Role{Name: roles[i].Name}).FirstOrCreate(&roles[i]).Error; err != nil {
			return err
		}
	}

	return nil
}
