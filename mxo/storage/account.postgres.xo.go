// Package storage contains the types for schema.
package storage

// Code generated by xo. DO NOT EDIT.

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// InsertAccount inserts the Account to the database.
func (s *PostgresStorage) InsertAccount(db XODB, a *Account) error {
	var err error

	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO "public"."account" (` +
		`"subject", "email", "created_date", "changed_date", "deleted_date"` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING "id"`

	// run query
	s.info(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate)
	err = db.QueryRow(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate).Scan(&a.ID)
	if err != nil {
		return err
	}

	// set existence
	a._exists = true

	return nil
}

// InsertAccountByFields inserts the Account to the database.
func (s *PostgresStorage) InsertAccountByFields(db XODB, a *Account) error {
	var err error

	params := make([]interface{}, 0, 5)
	fields := make([]string, 0, 5)
	retCols := `"id"`
	retVars := make([]interface{}, 0, 5)
	retVars = append(retVars, &a.ID)
	fields = append(fields, `"subject"`)
	params = append(params, a.Subject)

	fields = append(fields, `"email"`)
	params = append(params, a.Email)
	if a.CreatedDate.Valid {
		fields = append(fields, `"created_date"`)
		params = append(params, a.CreatedDate)
	} else {
		retCols += `, "created_date"`
		retVars = append(retVars, &a.CreatedDate)
	}
	if a.ChangedDate.Valid {
		fields = append(fields, `"changed_date"`)
		params = append(params, a.ChangedDate)
	} else {
		retCols += `, "changed_date"`
		retVars = append(retVars, &a.ChangedDate)
	}
	if a.DeletedDate.Valid {
		fields = append(fields, `"deleted_date"`)
		params = append(params, a.DeletedDate)
	} else {
		retCols += `, "deleted_date"`
		retVars = append(retVars, &a.DeletedDate)
	}
	if len(params) == 0 {
		// FIXME(jackie): maybe we should allow this?
		return errors.New("all fields are empty, unable to insert")
	}

	var placeHolders []string
	var placeHolderVals []interface{}
	for i := range params {
		placeHolders = append(placeHolders, "$%d")
		placeHolderVals = append(placeHolderVals, i+1)
	}
	placeHolderStr := fmt.Sprintf(strings.Join(placeHolders, ","), placeHolderVals...)

	sqlstr := `INSERT INTO "public"."account" (` +
		strings.Join(fields, ",") +
		`) VALUES (` + placeHolderStr +
		`) RETURNING ` + retCols

	s.info(sqlstr, params)
	err = db.QueryRow(sqlstr, params...).Scan(retVars...)
	if err != nil {
		return err
	}

	// set existence
	a._exists = true

	return nil
}

// UpdateAccount updates the Account in the database.
func (s *PostgresStorage) UpdateAccount(db XODB, a *Account) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if a._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query

	const sqlstr = `UPDATE "public"."account" SET (` +
		`"subject", "email", "created_date", "changed_date", "deleted_date"` +
		`) = ( ` +
		`$1, $2, $3, $4, $5` +
		`) WHERE "id" = $6`

	// run query
	s.info(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate, a.ID)
	_, err = db.Exec(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate, a.ID)
	return err
}

// UpdateAccountByFields updates the Account in the database.
func (s *PostgresStorage) UpdateAccountByFields(db XODB, a *Account, fields, retCols []string, params, retVars []interface{}) error {
	var placeHolders []string
	var idxvals []interface{}
	for i := range params {
		placeHolders = append(placeHolders, "$%d")
		idxvals = append(idxvals, i+1)
	}
	params = append(params, a.ID)
	idxvals = append(idxvals, len(params))

	var sqlstr string
	if len(fields) == 1 {
		sqlstr = fmt.Sprintf(`UPDATE "public"."account" SET `+
			strings.Join(fields, ",")+
			` = `+strings.Join(placeHolders, ",")+
			` WHERE id = $%d`+
			` RETURNING `+strings.Join(retCols, ", "), idxvals...)
	} else {
		sqlstr = fmt.Sprintf(`UPDATE "public"."account" SET (`+
			strings.Join(fields, ",")+
			`) = (`+strings.Join(placeHolders, ",")+
			`) WHERE id = $%d`+
			` RETURNING `+strings.Join(retCols, ", "), idxvals...)
	}
	s.info(sqlstr, params)
	if err := db.QueryRow(sqlstr, params...).Scan(retVars...); err != nil {
		return err
	}

	return nil
}

// SaveAccount saves the Account to the database.
func (s *PostgresStorage) SaveAccount(db XODB, a *Account) error {
	if a.Exists() {
		return s.UpdateAccount(db, a)
	}

	return s.InsertAccount(db, a)
}

// UpsertAccount performs an upsert for Account.
func (s *PostgresStorage) UpsertAccount(db XODB, a *Account) error {
	var err error

	// sql query
	const sqlstr = `INSERT INTO "public"."account" (` +
		`"id", "subject", "email", "created_date", "changed_date", "deleted_date"` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) ON CONFLICT ("id") DO UPDATE SET (` +
		`"id", "subject", "email", "created_date", "changed_date", "deleted_date"` +
		`) = (` +
		`EXCLUDED."id", EXCLUDED."subject", EXCLUDED."email", EXCLUDED."created_date", EXCLUDED."changed_date", EXCLUDED."deleted_date"` +
		`)`

	// run query
	s.info(sqlstr, a.ID, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate)
	_, err = db.Exec(sqlstr, a.ID, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate)
	if err != nil {
		return err
	}

	// set existence
	a._exists = true

	return nil
}

// DeleteAccount deletes the Account from the database.
func (s *PostgresStorage) DeleteAccount(db XODB, a *Account) error {
	var err error

	// if doesn't exist, bail
	if !a._exists {
		return nil
	}

	// if deleted, bail
	if a._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM "public"."account" WHERE "id" = $1`

	// run query
	s.info(sqlstr, a.ID)
	_, err = db.Exec(sqlstr, a.ID)
	if err != nil {
		return err
	}

	// set deleted
	a._deleted = true

	return nil
}

// DeleteAccounts deletes the Account from the database.
func (s *PostgresStorage) DeleteAccounts(db XODB, as []*Account) error {
	var err error

	if len(as) == 0 {
		return nil
	}

	var args []interface{}
	var placeholder string
	for i, a := range as {
		args = append(args, a.ID)
		if i != 0 {
			placeholder = placeholder + ", "
		}
		placeholder += fmt.Sprintf("$%d", i+1)
	}

	// sql query
	var sqlstr = `DELETE FROM "public"."account" WHERE "id" in (` + placeholder + `)`

	// run query
	s.info(sqlstr, args)
	_, err = db.Exec(sqlstr, args...)
	if err != nil {
		return err
	}

	// set deleted
	for _, a := range as {
		a._deleted = true
	}

	return nil
}

// GetMostRecentAccount returns n most recent rows from 'account',
// ordered by "created_date" in descending order.
func (s *PostgresStorage) GetMostRecentAccount(db XODB, n int) ([]*Account, error) {
	const sqlstr = `SELECT ` +
		`"id", "subject", "email", "created_date", "changed_date", "deleted_date" ` +
		`FROM "public"."account" ` +
		`ORDER BY created_date DESC LIMIT $1`

	s.info(sqlstr, n)
	q, err := db.Query(sqlstr, n)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Account
	for q.Next() {
		a := Account{}

		// scan
		err = q.Scan(&a.ID, &a.Subject, &a.Email, &a.CreatedDate, &a.ChangedDate, &a.DeletedDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
}

// GetMostRecentChangedAccount returns n most recent rows from 'account',
// ordered by "changed_date" in descending order.
func (s *PostgresStorage) GetMostRecentChangedAccount(db XODB, n int) ([]*Account, error) {
	const sqlstr = `SELECT ` +
		`"id", "subject", "email", "created_date", "changed_date", "deleted_date" ` +
		`FROM "public"."account" ` +
		`ORDER BY changed_date DESC LIMIT $1`

	s.info(sqlstr, n)
	q, err := db.Query(sqlstr, n)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Account
	for q.Next() {
		a := Account{}

		// scan
		err = q.Scan(&a.ID, &a.Subject, &a.Email, &a.CreatedDate, &a.ChangedDate, &a.DeletedDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
}

// GetAllAccount returns all rows from 'account', based on the AccountQueryArguments.
// If the AccountQueryArguments is nil, it will use the default AccountQueryArguments instead.
func (s *PostgresStorage) GetAllAccount(db XODB, queryArgs *AccountQueryArguments) ([]*Account, error) { // nolint: gocyclo
	queryArgs = ApplyAccountQueryArgsDefaults(queryArgs)
	if queryArgs.filterArgs == nil {
		filterArgs, err := getAccountFilter(queryArgs.Where)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get Account filter")
		}
		queryArgs.filterArgs = filterArgs
	}

	desc := ""
	if *queryArgs.Desc {
		desc = "DESC"
	}

	dead := "NULL"
	if *queryArgs.Dead {
		dead = "NOT NULL"
	}

	orderBy := "id"
	foundIndex := false
	dbFields := map[string]bool{
		"id":           true,
		"subject":      true,
		"email":        true,
		"created_date": true,
		"changed_date": true,
		"deleted_date": true,
	}

	if *queryArgs.OrderBy != "" && *queryArgs.OrderBy != defaultOrderBy {
		foundIndex = dbFields[*queryArgs.OrderBy]
		if !foundIndex {
			return nil, fmt.Errorf("unable to order by %s, field not found", *queryArgs.OrderBy)
		}
		orderBy = *queryArgs.OrderBy
	}

	var params []interface{}
	placeHolders := ""
	if queryArgs.filterArgs != nil {
		pls := make([]string, len(queryArgs.filterArgs.filterPairs))
		for i, pair := range queryArgs.filterArgs.filterPairs {
			pls[i] = fmt.Sprintf("%s %s $%d", pair.fieldName, pair.option, i+1)
			params = append(params, pair.value)
		}
		placeHolders = strings.Join(pls, " "+queryArgs.filterArgs.conjunction+" ")
		placeHolders = fmt.Sprintf("(%s) AND", placeHolders)
	}
	params = append(params, *queryArgs.Offset)
	offsetPos := len(params)

	params = append(params, *queryArgs.Limit)
	limitPos := len(params)

	var sqlstr = fmt.Sprintf(`SELECT %s FROM %s WHERE %s deleted_date IS %s ORDER BY %s %s OFFSET $%d LIMIT $%d`,
		`"id", "subject", "email", "created_date", "changed_date", "deleted_date" `,
		`"public"."account"`,
		placeHolders,
		dead,
		orderBy,
		desc,
		offsetPos,
		limitPos)
	s.info(sqlstr, params)

	q, err := db.Query(sqlstr, params...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Account
	for q.Next() {
		a := Account{}

		// scan
		err = q.Scan(&a.ID, &a.Subject, &a.Email, &a.CreatedDate, &a.ChangedDate, &a.DeletedDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
}

// CountAllAccount returns a count of all rows from 'account'
func (s *PostgresStorage) CountAllAccount(db XODB, queryArgs *AccountQueryArguments) (int, error) {
	queryArgs = ApplyAccountQueryArgsDefaults(queryArgs)
	if queryArgs.filterArgs == nil {
		filterArgs, err := getAccountFilter(queryArgs.Where)
		if err != nil {
			return 0, errors.Wrap(err, "unable to get Account filter")
		}
		queryArgs.filterArgs = filterArgs
	}

	dead := "NULL"
	if *queryArgs.Dead {
		dead = "NOT NULL"
	}

	var params []interface{}
	placeHolders := ""
	if queryArgs.filterArgs != nil {
		pls := make([]string, len(queryArgs.filterArgs.filterPairs))
		for i, pair := range queryArgs.filterArgs.filterPairs {
			pls[i] = fmt.Sprintf("%s %s $%d", pair.fieldName, pair.option, i+1)
			params = append(params, pair.value)
		}
		placeHolders = strings.Join(pls, " "+queryArgs.filterArgs.conjunction+" ")
		placeHolders = fmt.Sprintf("(%s) AND", placeHolders)
	}

	var err error
	var sqlstr = fmt.Sprintf(`SELECT count(*) from "public"."account" WHERE %s deleted_date IS %s`, placeHolders, dead)
	s.info(sqlstr)

	var count int
	err = db.QueryRow(sqlstr, params...).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// AccountByID retrieves a row from '"public"."account"' as a Account.
//
// Generated from index 'account_pk'.
func (s *PostgresStorage) AccountByID(db XODB, id int) (*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`"id", "subject", "email", "created_date", "changed_date", "deleted_date" ` +
		`FROM "public"."account" ` +
		`WHERE "id" = $1`

	// run query
	s.info(sqlstr, id)
	a := Account{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&a.ID, &a.Subject, &a.Email, &a.CreatedDate, &a.ChangedDate, &a.DeletedDate)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// AccountBySubject retrieves a row from '"public"."account"' as a Account.
//
// Generated from index 'account_subject_unique_index'.
func (s *PostgresStorage) AccountBySubject(db XODB, subject string) (*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`"id", "subject", "email", "created_date", "changed_date", "deleted_date" ` +
		`FROM "public"."account" ` +
		`WHERE "subject" = $1`

	// run query
	s.info(sqlstr, subject)
	a := Account{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, subject).Scan(&a.ID, &a.Subject, &a.Email, &a.CreatedDate, &a.ChangedDate, &a.DeletedDate)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
