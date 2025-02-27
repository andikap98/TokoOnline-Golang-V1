package model

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Product struct {
	ID       string `json:"id" binding:"len=0"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	IsDelete *bool  `json:"is_delete,omitempty"` //omitempty gunanya tidak menampilkan field kosong
}

var (
	ErrDBNil = errors.New("Koneksi Tidak Tersedia")
)

func SelectProduct(db *sql.DB) ([]Product, error) {
	if db == nil {
		return nil, ErrDBNil
	}
	query := `SELECT id, name, price FROM products WHERE is_delete = false`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil

}

func SelectProductIn(db *sql.DB, ids []string) ([]Product, error) {
	if db == nil {
		return nil, ErrDBNil
	}
	placeholder := []string{}
	arg := []any{}
	for i, id := range ids {
		placeholder = append(placeholder, fmt.Sprintf("$%d", i+1))
		arg = append(arg, id)
	}
	query := fmt.Sprintf(`SELECT id, name, price FROM products WHERE is_delete = false AND id IN (%s);`, strings.Join(placeholder, ","))
	rows, err := db.Query(query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func SelectProductByID(db *sql.DB, id string) (*Product, error) {
	if db == nil {
		return &Product{}, ErrDBNil
	}
	query := `SELECT id, name, price FROM products WHERE is_delete = false AND id = $1`
	row := db.QueryRow(query, id)

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func InsertProduct(db *sql.DB, product Product) error {

	if db == nil {
		return ErrDBNil
	}
	query := `INSERT INTO products (id, name, price) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProduct(db *sql.DB, product Product) error {
	if db == nil {
		return ErrDBNil
	}
	query := `UPDATE products SET name=$1, price=$2 WHERE id=$3`
	_, err := db.Exec(query, product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProduct(db *sql.DB, id string) error {
	if db == nil {
		return ErrDBNil
	}
	query := `UPDATE products SET is_delete=TRUE WHERE id=$1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
