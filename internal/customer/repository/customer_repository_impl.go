package repository

import (
	"fmt"
	"ps-eniqilo-store/internal/customer/dto"
	"ps-eniqilo-store/internal/customer/model"
	"ps-eniqilo-store/pkg/errs"
	"strings"

	"github.com/jmoiron/sqlx"
)

type customerRepository struct {
	db *sqlx.DB
}

func NewCustomerRepositoryImpl(db *sqlx.DB) CustomerRepository {
	return &customerRepository{db: db}
}

var (
	queryGetCustomerByEmail = `
    SELECT * FROM customers WHERE email = $1 LIMIT 1
`
)

func (r *customerRepository) CustomerEmailExists(email string) (bool, error) {
	var customer model.Customer
	err := r.db.Get(&customer, queryGetCustomerByEmail, email)
	if err != nil {
		return false, errs.NewErrInternalServerErrors("execute query error [GetCustomerByEmail]: ", err.Error())
	}

	if customer.ID != 0 {
		return true, nil
	}

	return false, err
}

var (
	queryGetCustomerByID = `
    SELECT * FROM customers WHERE id = $1 LIMIT 1
`
)

func (r *customerRepository) GetCustomerByID(customerId int64) (model.Customer, error) {
	var customer model.Customer
	err := r.db.Get(&customer, queryGetCustomerByID, customerId)
	if err != nil {
		return model.Customer{}, errs.NewErrInternalServerErrors("execute query error [GetCustomerByID]: ", err.Error())
	}
	return customer, err
}

var (
	queryCreateCustomer = `
    INSERT INTO customers (name, phone_number) 
    VALUES ($1, $2) 
    RETURNING id, name, phone_number
`
)

func (r *customerRepository) CreateCustomer(request *dto.CustomerReq) (model.Customer, error) {
	savedCustomer := model.Customer{}

	stmt, err := r.db.Prepare(queryCreateCustomer)
	if err != nil {
		return savedCustomer, errs.NewErrInternalServerErrors("query execute error on [CreateCustomer] : ", err.Error())
	}
	defer stmt.Close()

	// Execute the SQL statement to insert data
	var id int64
	var name string
	var phoneNumber string

	err = stmt.QueryRow(
		request.Name,
		request.PhoneNumber).Scan(&id, &name, &phoneNumber)

	if err != nil {
		if strings.Contains(err.Error(), "unique_phone_number") {
			return model.Customer{}, errs.NewErrDataConflict("phone number already exist", request.PhoneNumber)
		}

		return savedCustomer, errs.NewErrInternalServerErrors("query execute error on [CreateCustomer] : ", err.Error())
	}

	savedCustomer.ID = id
	savedCustomer.Name = name
	savedCustomer.PhoneNumber = phoneNumber

	return savedCustomer, nil
}

func (r *customerRepository) SearchCustomers(params map[string]interface{}) ([]model.Customer, error) {
	query := "SELECT id, name, phone_number FROM customers WHERE 1=1"
	args := []interface{}{}

	if name, ok := params["name"]; ok {
		query += " AND name LIKE $" + fmt.Sprintf("%d", len(args)+1)
		nameStr := fmt.Sprintf("%%%s%%", name)
		args = append(args, nameStr)
	}
  
  if phoneNumber, ok := params["phoneNumber"]; ok {
		query += " AND phone_number LIKE $" + fmt.Sprintf("%d", len(args)+1)
		phoneNumberStr := fmt.Sprintf("%%%s%%", phoneNumber)
		args = append(args, phoneNumberStr)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.PhoneNumber)
		if err != nil {
			return nil, errs.NewErrInternalServerErrors("execute query error [GetCustomer]: ", err.Error())
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.NewErrInternalServerErrors("execute query error [GetCustomer]: ", err.Error())
	}

	return customers, nil
}
