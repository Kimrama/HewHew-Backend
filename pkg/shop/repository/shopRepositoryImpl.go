package repository

import (
	"hewhew-backend/config"
	"hewhew-backend/database"
	"hewhew-backend/entities"
)
type ShopRepositoryImpl struct {
    db database.Database
    supabaseConfig *config.Supabase
}
func NewShopRepositoryImpl(db database.Database, supabaseConfig *config.Supabase) ShopRepository {
return &ShopRepositoryImpl{
    db: db,
    supabaseConfig: supabaseConfig,
}
}

func (r *ShopRepositoryImpl) CreateCanteen(canteenModel interface{}) error {
   return r.db.Connect().Create(canteenModel).Error
}   

func (r *ShopRepositoryImpl) EditCanteen(canteenName string,canteen *entities.Canteen) error {
	db := r.db.Connect()
	err := db.Model(&entities.Canteen{}).
		Where("canteen_name = ?", canteenName).
		Updates(map[string]interface{}{
			"latitude": canteen.Latitude,
			"longitude": canteen.Longitude,
		}).Error
	return err
}

func (r *ShopRepositoryImpl) DeleteCanteen(canteenName string) error {
	db := r.db.Connect()
	err := db.Where("canteen_name = ?", canteenName).Delete(&entities.Canteen{}).Error
	return err
	
}

