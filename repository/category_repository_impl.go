package repository

import (
	"belajar-golang-restfull-api/helper"
	"belajar-golang-restfull-api/model/domain"
	"context"
	"database/sql"
	"errors"
)

type CategoryRepositoryImpl struct {
}

func (c CategoryRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, Category domain.Category) domain.Category {
	SQL := "INSERT INTO category(name) VALUES (?);"
	result, err := tx.ExecContext(ctx, SQL, Category.Name)
	helper.PanicIfError(err)

	id, err := result.LastInsertId()
	helper.PanicIfError(err)

	Category.Id = int(id)
	return Category
}

func (c CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, Category domain.Category) domain.Category {
	SQL := "UPDATE category SET name = ? WHERE id = ?;"
	_, err := tx.ExecContext(ctx, SQL, Category.Name, Category.Id)
	helper.PanicIfError(err)

	return Category
}

func (c CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, Category domain.Category) {
	SQL := "DELETE FROM category WHERE id = ?;"
	_, err := tx.ExecContext(ctx, SQL, Category.Id)
	helper.PanicIfError(err)
}

func (c CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, CategoryId int) (domain.Category, error) {
	SQL := "SELECT * FROM category WHERE id = ?;"
	rows, err := tx.QueryContext(ctx, SQL, CategoryId)
	helper.PanicIfError(err)

	Category := domain.Category{}

	err = rows.Close()
	helper.PanicIfError(err)

	if rows.Next() {
		err = rows.Scan(&Category.Id, &Category.Name)
		helper.PanicIfError(err)
		return Category, nil
	} else {
		return Category, errors.New("category is not found")
	}
}

func (c CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	SQL := "SELECT * FROM category;"
	rows, err := tx.QueryContext(ctx, SQL)

	helper.PanicIfError(err)
	var categories []domain.Category

	for rows.Next() {
		category := domain.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err)
		categories = append(categories, category)
	}
	return categories
}
