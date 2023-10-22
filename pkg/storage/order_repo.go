package storage

import (
	"L0/pkg/model"
	"github.com/jackc/pgx"
	"log"
)

// OrderRepo - интерфейс репозитория для заказов
type OrderRepo interface {
	Insert(order model.Order) error
	GetById(uid string) (model.Order, bool)
}

// OrderRepoImpl - реализация репозитория для заказов. Содержит подключение к БД и cache заказов
type OrderRepoImpl struct {
	db    *pgx.Conn
	cache map[string]model.Order
}

func InitOrderRepo(db *pgx.Conn) *OrderRepoImpl {
	repo := &OrderRepoImpl{db: db, cache: map[string]model.Order{}}
	repo.loadCache()
	return repo
}

// Insert - метод для добавления заказа (в БД и кэш)
func (r *OrderRepoImpl) Insert(order model.Order) error {
	r.cache[order.Uid] = order
	_, err := r.db.Exec("INSERT INTO orders (uid, info) VALUES ($1, $2)", order.Uid, order)
	return err
}

// GetById - метод для получения заказа из кэша по id
func (r *OrderRepoImpl) GetById(uid string) (model.Order, bool) {
	order, ok := r.cache[uid]
	return order, ok
}

// loadCache - метод для загрузки заказов из БД в кэш
func (r *OrderRepoImpl) loadCache() {
	rows, err := r.db.Query("SELECT uid, info FROM orders")
	for rows.Next() {
		var uid string
		var order model.Order
		err = rows.Scan(&uid, &order)
		if err != nil {
			log.Fatal("Error to Scan SQL Orders in cache")
		}
		r.cache[uid] = order
	}
}
