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
	queryGetProductByID = `
    SELECT * FROM products WHERE id = $1
`
)

func (r *productRepository) GetProductByID(productId int64) (model.Product, error) {
	var product model.Product
	err := r.db.Get(&product, queryGetProductByID, productId)
	if err != nil {
		return model.Product{}, errs.NewErrInternalServerErrors("execute query error [GetProductByID]: ", err.Error())
	}
	return product, err
}

var (
	queryCreateProduct = `
    INSERT INTO products (name, sku, category, image_url, notes, price, stock, location, is_available) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
    RETURNING id, created_at
`
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

var (
	queryUpdateProduct = `
    UPDATE products
    SET name = $1, 
        sku = $2, 
        category = $3, 
        image_url = $4, 
        notes = $5, 
        price = $6, 
        stock = $7, 
        location = $8, 
        is_available = $9
    WHERE id = $10
    RETURNING created_at
`
)

func (r *productRepository) UpdateProduct(request *dto.ProductReq, productId int64) (model.Product, error) {
	savedProduct := model.Product{}

	stmt, err := r.db.Prepare(queryUpdateProduct)
	if err != nil {
		return savedProduct, errs.NewErrInternalServerErrors("query execute error on [UpdateProduct] : ", err.Error())
	}
	defer stmt.Close()

	// Execute the SQL statement to insert data
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
		request.IsAvailable,
		productId).Scan(&createdAt)
	if err != nil {
		return savedProduct, errs.NewErrInternalServerErrors("query execute error on [UpdateProduct] : ", err.Error())
	}
	savedProduct.ID = productId
	savedProduct.CreatedAtFormatter = createdAt.Format(time.RFC3339)
	return savedProduct, nil
}
