
package repository
import ("hewhew-backend/database")
type OrderRepositoryImpl struct {
    db database.Database
}
func NewOrderRepositoryImpl(db database.Database) OrderRepository {
return &OrderRepositoryImpl{
    db: db,
}
}
