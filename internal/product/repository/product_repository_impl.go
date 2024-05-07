package repository

import (
	"github.com/jmoiron/sqlx"
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/model"
	"ps-eniqilo-store/pkg/errs"
	"time"
)

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepositoryImpl(db *sqlx.DB) ProductRepository {
	return &productRepository{db: db}
}

var (
	queryCreateProduct = "INSERT INTO products (name, sku, category, image_url, notes, price, stock, location, is_available) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at"
)

func (r *productRepository) CreateProduct(request *dto.ProductReq) (model.Product, error) {
	savedProduct := model.Product{}

	stmt, err := r.db.Prepare(queryCreateProduct)
	if err != nil {
		return savedProduct, errs.NewErrInternalServerErrors("query execute error on [CreateProduct] : ", err.Error())
	}
	defer stmt.Close()

	// Execute the SQL statement to insert data
	var id int64
	var createdAt time.Time
	err = stmt.QueryRow(
		request.Name,
		request.SKU,
		request.Category,
		request.ImageURL,
		request.Notes,
		request.Price,
		request.Stock,
		request.Location,
		request.IsAvailable).Scan(&id, &createdAt)
	if err != nil {
		return savedProduct, errs.NewErrInternalServerErrors("query execute error on [CreateProduct] : ", err.Error())
	}
	savedProduct.ID = id
	savedProduct.CreatedAtFormatter = createdAt.Format(time.RFC3339)
	return savedProduct, nil
}
