package repository

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"ps-eniqilo-store/internal/checkout/dto"
	"ps-eniqilo-store/pkg/errs"
	"strconv"
	"strings"
	"time"
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
