package repository

import (
	"github.com/jmoiron/sqlx"
	"ps-eniqilo-store/internal/product/dto"
	"ps-eniqilo-store/internal/product/model"
	"ps-eniqilo-store/pkg/errs"
	"ps-eniqilo-store/pkg/helper"
	"strconv"
	"strings"
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
    SELECT * FROM products WHERE deleted_at IS NULL and id = $1
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

func generateQueryGetProducts(params map[string]interface{}) (string, []interface{}) {
	query := "SELECT * FROM products WHERE deleted_at IS NULL and 1=1"
	var orderByParts []string
	isOrder := false

	var args []interface{}
	num := 1
	limit := 5
	offset := 0
	for key, value := range params {
		isAddArgs := false
		switch key {
		case "id":
			query += " AND id = $" + strconv.Itoa(num)
			isAddArgs = true
			num++
		case "name":
			query += " AND name LIKE '%' || $" + strconv.Itoa(num) + " || '%'"
			isAddArgs = true
			num++
		case "category":
			query += " AND category = $" + strconv.Itoa(num)
			isAddArgs = true
			num++
		case "isAvailable":
			query += " AND is_available = $" + strconv.Itoa(num)
			isAddArgs = true
			num++
		case "sku":
			query += " AND sku = $" + strconv.Itoa(num)
			isAddArgs = true
			num++
		case "inStock":
			if value.(bool) {
				query += " AND stock > 0"
			} else {
				query += " AND stock = 0"
			}
		case "price":
			orderByParts = append(orderByParts, " price "+value.(string))
			isOrder = true
		case "createdAt":
			orderByParts = append(orderByParts, " created_at "+value.(string))
			isOrder = true
		case "limit":
			limit = value.(int)
		case "offset":
			offset = value.(int)
		}
		if isAddArgs {
			args = append(args, value)
		}
	}

	if isOrder {
		query += " ORDER BY " + strings.Join(orderByParts, ", ")
	}
	query += " LIMIT $" + strconv.Itoa(num) + " OFFSET $" + strconv.Itoa(num+1)
	args = append(args, limit)
	args = append(args, offset)
	return query, args
}

func (r *productRepository) GetProducts(params map[string]interface{}) ([]model.Product, error) {
	query, args := generateQueryGetProducts(params)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Category,
			&product.ImageURL,
			&product.Notes,
			&product.Price,
			&product.Stock,
			&product.Location,
			&product.IsAvailable,
			&product.CreatedAt,
			&product.DeletedAt)
		if err != nil {
			return nil, errs.NewErrInternalServerErrors("execute query error [GetProducts]: ", err.Error())
		}
		product.IDString = helper.IntToString(product.ID)
		product.CreatedAtFormatter = product.CreatedAt.Format(time.RFC3339)
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, errs.NewErrInternalServerErrors("execute query error [GetProducts]: ", err.Error())
	}

	return products, nil
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

var (
	queryDeleteProduct = `
    UPDATE products SET deleted_at = NOW()  WHERE id = $1
`
)

func (r *productRepository) DeleteProduct(productId int64) error {
	stmt, err := r.db.Prepare(queryDeleteProduct)
	if err != nil {
		return errs.NewErrInternalServerErrors("query execute error on [DeleteProduct] : ", err.Error())
	}
	defer stmt.Close()

	// Execute the SQL statement to insert data
	_, err = stmt.Query(productId)
	if err != nil {
		return errs.NewErrInternalServerErrors("query execute error on [DeleteProduct] : ", err.Error())
	}

	return nil
}
