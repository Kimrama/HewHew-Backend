
package repository
import ("hewhew-backend/database")
type ShopRepositoryImpl struct {
    db database.Database
}
func NewShopRepositoryImpl(db database.Database) ShopRepository {
return &ShopRepositoryImpl{
    db: db,
}
}
