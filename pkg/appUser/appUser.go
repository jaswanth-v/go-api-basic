// User type and associated methods to create, modify, delete application users
package appUser

import (
	"context"
	"errors"

	"github.com/gilcrest/go-API-template/pkg/config/db"
)

type User struct {
	Username  string
	MobileID  string
	Email     string
	FirstName string
	LastName  string
}

// Perform business validations prior to writing to the db
func (u User) Create(ctx context.Context) (int, error) {

	if u.Email == "" {
		return -1, errors.New("Email must have a value")
	}

	// Write to db -- createDB method returns rows impacted (should always be 1)
	// or an error
	rows, err := u.createDB(ctx)
	return rows, err
}

// Creates a record in the appUser table using a stored function which
// returns the number of rows inserted
func (u User) createDB(ctx context.Context) (int, error) {

	var (
		rowsInserted int
	)

	// pull pointer to sql.Tx as tx from context passed in parameter
	tx, ok := db.DBTxFromContext(ctx)

	// ensure there is a sql.Tx in the context by checking boolean passed back above
	if !ok {
		return -1, errors.New("Unable to retrieve sql.Tx from DBTxFromContext function")
	}

	// Prepare the sql statement using bind variables
	stmt, err := tx.PrepareContext(ctx, "select lp.create_app_user(p_username => $1, p_mobile_id => $2, p_email_address => $3, p_first_name => $4, p_last_name => $5, p_create_user_id => $6)")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	// Execute stored function that returns rows impacted, hence the use of QueryContext instead of Exec
	rows, err := stmt.QueryContext(ctx, u.Username, u.MobileID, u.Email, u.FirstName, u.LastName, "gilcrest")
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	// Iterate through
	for rows.Next() {
		if err := rows.Scan(&rowsInserted); err != nil {
			return -1, err
		}
	}

	if err := rows.Err(); err != nil {
		return -1, err
	}

	return rowsInserted, nil

}