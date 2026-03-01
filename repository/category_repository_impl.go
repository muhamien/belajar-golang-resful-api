package repository

import (
	"belajar-golang-restful-api/helper"
	"belajar-golang-restful-api/model/domain"
	"context"
	"database/sql"
	"errors"
)

type CategoryRepositoryImpl struct {
}

func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{}
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

func (c CategoryRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, category domain.Category) domain.Category {
	SQL := "UPDATE category SET name = ? WHERE id = ?;"
	_, err := tx.ExecContext(ctx, SQL, category.Name, category.Id)
	helper.PanicIfError(err)

	return category
}

func (c CategoryRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, Category domain.Category) {
	SQL := "DELETE FROM category WHERE id = ?;"
	_, err := tx.ExecContext(ctx, SQL, Category.Id)
	helper.PanicIfError(err)
}

func (c CategoryRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, CategoryId int) (domain.Category, error) {
	SQL := "SELECT id, name FROM category WHERE id = ?;"
	rows, err := tx.QueryContext(ctx, SQL, CategoryId)
	helper.PanicIfError(err)
	defer rows.Close()

	Category := domain.Category{}

	if rows.Next() {
		err = rows.Scan(&Category.Id, &Category.Name)
		helper.PanicIfError(err)
		return Category, nil
	} else {
		return Category, errors.New("category is not found")
	}
}

func (c CategoryRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) []domain.Category {
	SQL := "SELECT id, name FROM category;"
	rows, err := tx.QueryContext(ctx, SQL)
	helper.PanicIfError(err)
	defer rows.Close()

	categories := []domain.Category{}

	for rows.Next() {
		category := domain.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		helper.PanicIfError(err)
		categories = append(categories, category)
	}
	return categories
}
