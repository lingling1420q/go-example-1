// Package storage contains the types for schema.
package storage

// Code generated by xo. DO NOT EDIT.

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// InsertAccount inserts the Account to the database.
func (s *MssqlStorage) InsertAccount(db XODB, a *Account) error {
	var err error

	// if already exist, bail
	if a._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by identity
	const sqlstr = `INSERT INTO dbo.account (` +
		`subject, email, created_date, changed_date, deleted_date` +
		`) OUTPUT INSERTED.ID VALUES (` +
		`$1, $2, $3, $4, $5` +
		`)`

	// run query
	s.info(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate)
	err = db.QueryRow(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate).Scan(&a.ID)
	if err != nil {
		return err
	}

	// set primary key and existence
	a._exists = true

	return nil
}

// InsertAccountByFields inserts the Account to the database.
func (s *MssqlStorage) InsertAccountByFields(db XODB, a *Account) error {
	var err error

	params := make([]interface{}, 0, 5)
	fields := make([]string, 0, 5)
	retCols := `INSERTED.id`
	retVars := make([]interface{}, 0, 5)
	retVars = append(retVars, &a.ID)
	fields = append(fields, `subject`)
	params = append(params, a.Subject)

	fields = append(fields, `email`)
	params = append(params, a.Email)
	if a.CreatedDate.Valid {
		fields = append(fields, `created_date`)
		params = append(params, a.CreatedDate)
	} else {
		retCols += `, INSERTED.created_date`
		retVars = append(retVars, &a.CreatedDate)
	}
	if a.ChangedDate.Valid {
		fields = append(fields, `changed_date`)
		params = append(params, a.ChangedDate)
	} else {
		retCols += `, INSERTED.changed_date`
		retVars = append(retVars, &a.ChangedDate)
	}
	if a.DeletedDate.Valid {
		fields = append(fields, `deleted_date`)
		params = append(params, a.DeletedDate)
	} else {
		retCols += `, INSERTED.deleted_date`
		retVars = append(retVars, &a.DeletedDate)
	}
	if len(params) == 0 {
		// FIXME(jackie): maybe we should allow this?
		return errors.New("all fields are empty, unable to insert")
	}

	var placeHolders string
	for i := range params {
		placeHolders += "$" + strconv.Itoa(i+1)
		if i < len(params)-1 {
			placeHolders += ", "
		}
	}

	sqlstr := `INSERT INTO dbo.account (` +
		strings.Join(fields, ",") +
		`) OUTPUT ` + retCols +
		` VALUES (` + placeHolders + `)`

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
func (s *MssqlStorage) UpdateAccount(db XODB, a *Account) error {
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
	const sqlstr = `UPDATE dbo.account SET ` +
		`subject = $1, email = $2, created_date = $3, changed_date = $4, deleted_date = $5` +
		` WHERE id = $6`

	// run query
	s.info(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate, a.ID)
	_, err = db.Exec(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate, a.ID)
	return err
}

// UpdateAccountByFields updates the Account in the database.
func (s *MssqlStorage) UpdateAccountByFields(db XODB, a *Account, fields, retCols []string, params, retVars []interface{}) error {
	var setstr string
	for i, field := range fields {
		if i != 0 {
			setstr += ", "
		}
		setstr += field + ` = $` + strconv.Itoa(i+1)
	}

	var retstr string
	for i, retCol := range retCols {
		if i != 0 {
			retstr += ", "
		}
		retstr += "INSERTED." + retCol
	}

	params = append(params, a.ID)
	var sqlstr = `UPDATE dbo.account SET ` +
		setstr + ` OUTPUT ` + retstr +
		` WHERE id = $` + strconv.Itoa(len(params))
	s.info(sqlstr, params)
	if err := db.QueryRow(sqlstr, params...).Scan(retVars...); err != nil {
		return err
	}

	return nil
}

// SaveAccount saves the Account to the database.
func (s *MssqlStorage) SaveAccount(db XODB, a *Account) error {
	if a.Exists() {
		return s.UpdateAccount(db, a)
	}

	return s.InsertAccount(db, a)
}

// UpsertAccount performs an upsert for Account.
func (s *MssqlStorage) UpsertAccount(db XODB, a *Account) error {
	var err error

	// sql query

	const sqlstr = `MERGE dbo.account AS t ` +
		`USING (SELECT $1 AS id, $2 AS subject, $3 AS email, $4 AS created_date, $5 AS changed_date, $6 AS deleted_date) AS s ` +
		`ON t.id = s.id ` +
		`WHEN MATCHED THEN UPDATE SET subject = s.subject, email = s.email, created_date = s.created_date, changed_date = s.changed_date, deleted_date = s.deleted_date ` +
		`WHEN NOT MATCHED THEN INSERT (subject, email, created_date, changed_date, deleted_date) VALUES (s.subject, s.email, s.created_date, s.changed_date, s.deleted_date);`

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
func (s *MssqlStorage) DeleteAccount(db XODB, a *Account) error {
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
	const sqlstr = `DELETE FROM dbo.account WHERE id = $1`

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
func (s *MssqlStorage) DeleteAccounts(db XODB, as []*Account) error {
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
	var sqlstr = `DELETE FROM dbo.account WHERE id in (` + placeholder + `)`

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
func (s *MssqlStorage) GetMostRecentAccount(db XODB, n int) ([]*Account, error) {
	var sqlstr = `SELECT TOP ` + strconv.Itoa(n) +
		` id, subject, email, created_date, changed_date, deleted_date ` +
		`FROM dbo.account ` +
		`ORDER BY created_date DESC`

	s.info(sqlstr)
	q, err := db.Query(sqlstr)
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
func (s *MssqlStorage) GetMostRecentChangedAccount(db XODB, n int) ([]*Account, error) {
	var sqlstr = `SELECT TOP ` + strconv.Itoa(n) +
		` id, subject, email, created_date, changed_date, deleted_date ` +
		`FROM dbo.account ` +
		`ORDER BY changed_date DESC`

	s.info(sqlstr)
	q, err := db.Query(sqlstr)
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
func (s *MssqlStorage) GetAllAccount(db XODB, queryArgs *AccountQueryArguments) ([]*Account, error) { // nolint: gocyclo
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

	var sqlstr = fmt.Sprintf(`SELECT %s FROM %s WHERE %s deleted_date IS %s ORDER BY %s %s OFFSET $%d ROWS FETCH NEXT $%d ROWS ONLY`,
		`id, subject, email, created_date, changed_date, deleted_date `,
		`dbo.account`,
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
func (s *MssqlStorage) CountAllAccount(db XODB, queryArgs *AccountQueryArguments) (int, error) {
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
	var sqlstr = fmt.Sprintf(`SELECT count(*) from dbo.account WHERE %s deleted_date IS %s`, placeHolders, dead)
	s.info(sqlstr)

	var count int
	err = db.QueryRow(sqlstr, params...).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// AccountByID retrieves a row from 'dbo.account' as a Account.
//
// Generated from index 'PK__account__3213E83F136498D2'.
func (s *MssqlStorage) AccountByID(db XODB, id int) (*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, subject, email, created_date, changed_date, deleted_date ` +
		`FROM dbo.account ` +
		`WHERE id = $1`

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

// AccountBySubject retrieves a row from 'dbo.account' as a Account.
//
// Generated from index 'account_subject_ak'.
func (s *MssqlStorage) AccountBySubject(db XODB, subject string) (*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, subject, email, created_date, changed_date, deleted_date ` +
		`FROM dbo.account ` +
		`WHERE subject = $1`

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
