package migration

import (
	"fmt"
	"hewhew-backend/database"
	"hewhew-backend/entities"

	"gorm.io/gorm"
)

func Migrate(db database.Database) {
	tx := db.Connect().Begin()

	// base entity (No FK)
	UserMigration(tx)
	ShopMigration(tx)
	WasteDropMigration(tx)
	// ==========

	// entities relate with base by FK
	ContactMigration(tx)
	TopUpMigration(tx)
	ShopAdminMigration(tx)
	MenuMigration(tx)
	// ==========

	// entities that need user/shop/menu
	OrderMigration(tx)
	FavouriteMigration(tx)
	ReviewMigration(tx)

	// ===========
	// entities that need order/menu
	MenuQuantityMigration(tx)
	ChatMigration(tx)
	NotificationMigration(tx)
	TransactionLogMigration(tx)
	// ===========

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		panic(err)
	}

	fmt.Println("Database migration completed.")
}

func UserMigration(tx *gorm.DB) error {
	fmt.Println("Migrating User table...")
	return tx.Migrator().CreateTable(&entities.User{})
}
func ShopMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Shop table...")
	return tx.Migrator().CreateTable(&entities.Shop{})
}
func WasteDropMigration(tx *gorm.DB) error {
	fmt.Println("Migrating WasteDrop table...")
	return tx.Migrator().CreateTable(&entities.WasteDrop{})
}
func ContactMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Contact table...")
	return tx.Migrator().CreateTable(&entities.Contact{})
}
func TopUpMigration(tx *gorm.DB) error {
	fmt.Println("Migrating TopUp table...")
	return tx.Migrator().CreateTable(&entities.TopUp{})
}
func ShopAdminMigration(tx *gorm.DB) error {
	fmt.Println("Migrating ShopAdmin table...")
	return tx.Migrator().CreateTable(&entities.ShopAdmin{})
}
func MenuMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Menu table...")
	return tx.Migrator().CreateTable(&entities.Menu{})
}
func OrderMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Order table...")
	return tx.Migrator().CreateTable(&entities.Order{})
}
func FavouriteMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Favourite table...")
	return tx.Migrator().CreateTable(&entities.Favourite{})
}
func ReviewMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Review table...")
	return tx.Migrator().CreateTable(&entities.Review{})
}
func MenuQuantityMigration(tx *gorm.DB) error {
	fmt.Println("Migrating MenuQuantity table...")
	return tx.Migrator().CreateTable(&entities.MenuQuantity{})
}
func ChatMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Chat table...")
	return tx.Migrator().CreateTable(&entities.Chat{})
}
func NotificationMigration(tx *gorm.DB) error {
	fmt.Println("Migrating Notification table...")
	return tx.Migrator().CreateTable(&entities.Notification{})
}
func TransactionLogMigration(tx *gorm.DB) error {
	fmt.Println("Migrating TransactionLog table...")
	return tx.Migrator().CreateTable(&entities.TransactionLog{})
}
