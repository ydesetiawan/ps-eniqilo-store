package repository

import (
	"encoding/json"
	"ps-eniqilo-store/internal/checkout/dto"
	"ps-eniqilo-store/internal/checkout/model"
	"ps-eniqilo-store/pkg/errs"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type checkoutRepository struct {
	db *sqlx.DB
}

func NewCheckoutRepositoryImpl(db *sqlx.DB) CheckoutRepository {
	return &checkoutRepository{db: db}
}

var (
	queryGetHistory = `
        SELECT 
        	CAST(c.id AS text) AS transactionId,
            CAST(c.customer_id AS text) AS customerId,
            JSON_AGG(jsonb_build_object(
                'productId', CAST(cd.product_id AS text),
                'quantity', cd.quantity
            )) AS productDetails,
            c.paid,  
            c.change,
            c.created_at
        FROM 
            checkouts c
        LEFT JOIN LATERAL (
            SELECT 
                product_id,
                quantity
            FROM 
                checkout_details cd
            WHERE 
                c.id = cd.checkout_id
        ) cd ON true
        WHERE 1=1
    `
)

func generateQueryGetHistory(params map[string]interface{}) (string, []interface{}) {
	query := queryGetHistory
	var orderByParts []string
	isOrder := false

	var args []interface{}
	num := 1
	limit := 5
	offset := 0
	for key, value := range params {
		isAddArgs := false
		switch key {
		case "customerId":
			query += " AND c.customer_id = $" + strconv.Itoa(num)
			isAddArgs = true
			num++
		case "createdAt":
			orderByParts = append(orderByParts, " c.created_at "+value.(string))
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

	query += " GROUP BY c.id, c.customer_id, c.paid, c.change"

	if isOrder {
		query += " ORDER BY " + strings.Join(orderByParts, ", ")
	}

	query += " LIMIT $" + strconv.Itoa(num) + " OFFSET $" + strconv.Itoa(num+1)
	args = append(args, limit)
	args = append(args, offset)
	return query, args
}

func (c checkoutRepository) GetCheckoutHistory(params map[string]interface{}) ([]dto.CheckOutHistoryResp, error) {

	query, args := generateQueryGetHistory(params)
	rows, err := c.db.Query(query, args...)
	if err != nil {
		return nil, errs.NewErrInternalServerErrors("execute query error [GetCheckoutHistory]: ", err.Error())
	}
	defer rows.Close()

	var histories []dto.CheckOutHistoryResp
	for rows.Next() {
		var resp dto.CheckOutHistoryResp
		var productDetailsJSON string
		var createdAtTime time.Time
		err := rows.Scan(&resp.TransactionId, &resp.CustomerId, &productDetailsJSON, &resp.Paid, &resp.Change, &createdAtTime)
		if err != nil {
			return nil, errs.NewErrInternalServerErrors("execute query error [GetCheckoutHistory]: ", err.Error())
		}

		err = json.Unmarshal([]byte(productDetailsJSON), &resp.ProductDetails)
		if err != nil {
			return nil, errs.NewErrInternalServerErrors("execute query error [GetCheckoutHistory]: ", err.Error())
		}

		resp.CreatedAt = createdAtTime.Format(time.RFC3339)
		histories = append(histories, resp)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.NewErrInternalServerErrors("execute query error [GetCheckoutHistory]: ", err.Error())
	}
	return histories, nil

}

const createCheckoutQuery = "INSERT INTO checkouts (customer_id, total_price, paid, change) VALUES ($1, $2, $3, $4) RETURNING id"

func (c checkoutRepository) CreateCheckout(tx *sqlx.Tx, checkout *model.Checkout) (checkoutId int64, err error) {
	if tx != nil {
		err = tx.QueryRowx(createCheckoutQuery, checkout.CustomerID, checkout.TotalPrice, checkout.Paid, checkout.Change).Scan(&checkoutId)
	} else {
		err = c.db.QueryRowx(createCheckoutQuery, checkout.CustomerID, checkout.TotalPrice, checkout.Paid, checkout.Change).Scan(&checkoutId)
	}

	return
}

const createCheckoutDetailQuery = "INSERT INTO checkout_details (checkout_id, product_id, product_price, total_price, quantity) VALUES ($1, $2, $3, $4, $5)"

func (c checkoutRepository) CreateCheckoutDetail(tx *sqlx.Tx, detail *model.CheckoutDetail) (err error) {
	if tx != nil {
		_, err = tx.Exec(createCheckoutDetailQuery, detail.CheckoutID, detail.ProductID, detail.ProductPrice, detail.TotalPrice, detail.Quantity)
	} else {
		_, err = c.db.Exec(createCheckoutDetailQuery, detail.CheckoutID, detail.ProductID, detail.ProductPrice, detail.TotalPrice, detail.Quantity)
	}
	return
}
