package database

import (
	"fmt"

	"gorm.io/gorm"

	"tienda/backend/internal/models"
)

// Migrate runs GORM's AutoMigrate on every model in the project. It
// creates missing tables, adds missing columns, and creates missing
// indexes/unique constraints declared in the struct tags. It does NOT
// drop or rename existing schema objects.
//
// Scope notes:
//   - Only the column/index side of each table is reconciled. Foreign-key
//     REFERENCES and ON DELETE policies declared in DBML live in
//     internal/database/database.sql (the source of truth for FK
//     topology). They are NOT created here — apply that file separately
//     before this migration in environments where the schema must match
//     the DBML exactly.
//   - Postgres CHECK constraints from DBML (e.g. chk_review_rating_range)
//     are also not created by AutoMigrate; they live with the FKs in
//     the DBML-generated SQL.
//
// Model order follows dependency direction (root tables first, joins last)
// for readability. AutoMigrate resolves the actual creation order from
// declared FK relationships, but a stable source order makes diffs in
// generated SQL easier to read.
func Migrate(db *gorm.DB) error {
	all := []any{
		// Root tables (no project FKs).
		&models.User{},
		&models.Admins{},
		&models.Permission{},
		&models.Category{},
		&models.Impresora3d{},
		&models.Filament{},

		// Items & printing configs (depend on root tables).
		&models.Item3D{},
		&models.Item{},

		// Cart & addresses (depend on User).
		&models.Cart{},
		&models.ShippingAddress{},

		// Cart line.
		&models.CartItem{},

		// Orders & payments.
		&models.Order{},
		&models.OrderItem{},
		&models.Payment{},

		// Item audit & reviews.
		&models.StockMovement{},
		&models.Review{},

		// Admin permission join.
		&models.AdminPermission{},
	}

	if err := db.AutoMigrate(all...); err != nil {
		return fmt.Errorf("database: automigrate: %w", err)
	}
	return nil
}
