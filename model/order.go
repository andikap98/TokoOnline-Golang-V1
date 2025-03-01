package model

import (
	"database/sql"
	"time"
)

type Checkout struct {
	Email    string            `json:"email"`
	Address  string            `json:"address"`
	Products []ProductQuantity `json:"products"`
}

type ProductQuantity struct {
	ID       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

type Order struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Address     string     `json:"address"`
	GrandTotal  int64      `json:"grandTotal"`
	Passcode    *string    `json:"passcode,omitempty"` // omitempty mengabaikan field dengan nilai kosong
	PaidAt      *time.Time `json:"paid_at,omitempty"`
	PaidBank    *string    `json:"paid_bank,omitempty"`
	PaidAccount *string    `json:"paid_account,omitempty"`
}

type OrderDetail struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	Price     int64  `json:"price"`
	Total     int64  `json:"total"`
}

type OrderWithDetail struct {
	Order
	Details []OrderDetail `json:"details"`
}

type Confirm struct {
	Amount        int64  `json:"amount" binding:"required"`
	Bank          string `json:"bank" binding:"required"`
	AccountNumber string `json:"accountNumber" binding:"required"`
	Passcode      string `json:"passcode" binding:"required"`
}

func CreateOrder(db *sql.DB, order Order, details []OrderDetail) error {
	if db == nil {
		return ErrDBNil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	queryOrder := `INSERT INTO orders (id, email, address, passcode, grand_total) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(queryOrder, order.ID, order.Email, order.Address, order.Passcode, order.GrandTotal)
	if err != nil {
		tx.Rollback()
		return err
	}
	queryDetails := `INSERT INTO order_details (id, order_id, product_id, quantity, price, total) VALUES ($1, $2, $3, $4, $5, $6);`
	for _, d := range details {
		_, err := tx.Exec(queryDetails, d.ID, d.OrderID, d.ProductID, d.Quantity, d.Price, d.Total)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func SelectOrderByID(db *sql.DB, id string) (Order, error) {
	if db == nil {
		return Order{}, ErrDBNil
	}

	query := `SELECT id, email, address, passcode, grand_total, paid_at, paid_bank, paid_account FROM orders WHERE id = $1`
	row := db.QueryRow(query, id)
	var order Order
	err := row.Scan(&order.ID, &order.Email, &order.Address, &order.Passcode, &order.GrandTotal, &order.PaidAt, &order.PaidBank, &order.PaidAccount)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func UpdateOrderByID(db *sql.DB, id string, confirm Confirm, paidAt time.Time) error {
	if db == nil {
		return ErrDBNil
	}
	query := `UPDATE orders SET paid_at = $1, paid_bank =$2, paid_account =$3 WHERE id = $4`
	if _, err := db.Exec(query, paidAt, confirm.Bank, confirm.AccountNumber, id); err != nil {
		return err
	}

	return nil

}
