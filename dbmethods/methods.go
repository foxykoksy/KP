package dbmethods
import (
	"database/sql"
	"fmt"
	"awesomeProject3/ProductsInfo"
	"log"
	"strconv"
)

type dbmethods interface {
	AddItem(p *ProductsInfo.Product) error
	GetItem(item string, company string) ([]ProductsInfo.Product, error)
	DeleteItem(item string, company string, price int, amount int) error
	AddInfo(p *ProductsInfo.Info) bool
	GetInfo(company string) (string, int, error)
	ShowAll() error
}

var count int

type dbm struct {
	db *sql.DB
}

func (d dbm) AddItem(p *ProductsInfo.Product) error {
	row := d.db.QueryRow("SELECT * FROM products WHERE item = $1 AND company = $2 AND price = $3 AND amount = $4;",
		p.Item, p.Company, p.Price, p.Amount)
	prod := ProductsInfo.Product{}
	err := row.Scan(&prod.Id, &prod.Item, &prod.Company, &prod.Price, &prod.Amount)
	if count == 0 || err == sql.ErrNoRows {
		_, err := d.db.Exec("INSERT INTO products (item, company, price, amount) VALUES ($1, $2, $3, $4);",
			p.Item, p.Company, p.Price, p.Amount)
		count++
		return err
	}
	_, err = d.db.Exec("UPDATE products SET amount = $1 WHERE id = $2;", prod.Amount+p.Amount, prod.Id)
	count++
	return err
}

func (d dbm) GetItem(item string, company string) ([]ProductsInfo.Product, error) {
	rows, err := d.db.Query("SELECT * FROM products WHERE item = $1 AND company = $2;", item, company)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	p := []ProductsInfo.Product{}
	for rows.Next() {
		prod := ProductsInfo.Product{}
		err := rows.Scan(&prod.Id, &prod.Item, &prod.Company, &prod.Price, &prod.Amount)
		if err != nil {
			fmt.Println(err)
			continue
		}
		p = append(p, prod)
	}
	return p, nil
}

func (d dbm) DeleteItem(item string, company string, amount int) error {
	if count != 0 {
		row := d.db.QueryRow("SELECT * FROM products WHERE item = $1 AND company = $2;", item, company)
		if row != nil {
			prod := ProductsInfo.Product{}
			err := row.Scan(&prod.Id, &prod.Item, &prod.Company, &prod.Price, &prod.Amount)
			if prod.Amount > amount {
				_, err = d.db.Exec("UPDATE products SET amount = $1 WHERE id = $2;", prod.Amount-amount, prod.Id)
			} else {
				_, err = d.db.Exec("DELETE FROM products WHERE item = $1 AND company = $2;", item, company)
			}
			return err
		}
		return nil
	}
	return nil
}

func (d dbm) ShowAll() error {
	rows, err := d.db.Query("SELECT * FROM products;")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		prod := ProductsInfo.Product{}
		_ = rows.Scan(&prod.Id, &prod.Item, &prod.Company, &prod.Price, &prod.Amount)
		ProductsInfo.ShowItem(prod)
	}
	return nil
}

func (d dbm) AddInfo(p *ProductsInfo.Info) bool {
	res := d.db.QueryRow("SELECT * FROM info WHERE company = $1;", p.Company)
	err := res.Scan(&p.Id, &p.Company, &p.Information, &p.Rating)
	if err != nil {
		_, err = d.db.Exec("INSERT INTO info (company, information, rating) VALUES ($1, $2, $3);",
			p.Company, p.Information, p.Rating)
		return true
	}
	return false
}

func (d dbm) GetInfo(company string) (string, error) {
	res := d.db.QueryRow("SELECT * FROM info WHERE company = $1;", company)
	prod := ProductsInfo.Info{}
	if res != nil {
		err := res.Scan(&prod.Id, &prod.Company, &prod.Information, &prod.Rating)
		s := prod.Information + "; Rating: " + strconv.Itoa(prod.Rating)
		return s, err
	}
	return "No information", nil
}

func NewItemTable() (*dbm, error) {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE if NOT EXISTS products (id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
		item TEXT NOT NULL,
		company TEXT NOT NULL,
		price INTEGER NOT NULL, 
		amount INTEGER DEFAULT 1)`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE if NOT EXISTS info (id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
		company TEXT NOT NULL,
		information TEXT NOT NULL,
		rating INTEGER NOT NULL)`)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return &dbm{db: db}, nil
}