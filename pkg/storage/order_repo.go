package storage

import (
	"L0/pkg/model"
	"github.com/jackc/pgx"
	"log"
)

type OrderRepo interface {
	Insert(order model.Order) error
	GetById(uid string) (model.Order, bool)
}

type OrderRepoImpl struct {
	db    *pgx.Conn
	cache map[string]model.Order
}

func InitOrderRepo(db *pgx.Conn) *OrderRepoImpl {
	repo := &OrderRepoImpl{db: db, cache: map[string]model.Order{}}
	repo.loadCache()
	return repo
}

func (r *OrderRepoImpl) Insert(order model.Order) error {
	r.cache[order.Uid] = order
	_, err := r.db.Exec("INSERT INTO orders (uid, info) VALUES ($1, $2)", order.Uid, order)
	return err
}

func (r *OrderRepoImpl) GetById(uid string) (model.Order, bool) {
	order, ok := r.cache[uid]
	return order, ok
}

func (r *OrderRepoImpl) loadCache() {
	rows, err := r.db.Query("SELECT uid, info FROM orders")
	for rows.Next() {
		var uid string
		var order model.Order
		err = rows.Scan(&uid, &order)
		if err != nil {
			log.Fatal("Error to Scan SQL Result")
		}
		r.cache[uid] = order
	}
}
