package app

import (
	"belajar-golang-restful-api/helper"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/belajar_golang_restful_api")
	helper.PanicIfError(err)

	// setup database pooling
	db.SetMaxIdleConns(5)                   // minimal jumlah koneksi standby (idle)
	db.SetMaxOpenConns(20)                  // maksimal jumlah koneksi yang bisa dibuka
	db.SetConnMaxLifetime(60 * time.Minute) // berapa lama koneksi boleh digunakan sebelum direfresh
	db.SetConnMaxIdleTime(10 * time.Minute) // berapa lama koneksi idle boleh bertahan sebelum dihapus

	return db
}
