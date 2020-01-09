// Package storage contains the types for schema.
package storage

// Code generated by xo. DO NOT EDIT.

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// InsertUser inserts the User to the database.
func (s *PostgresStorage) InsertUser(db XODB, u *User) error {
	var err error

	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by sequence
	const sqlstr = `INSERT INTO "public"."user" (` +
		`"subject", "name", "created_date", "changed_date", "deleted_date"` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING "id"`

	// run query
	XOLog(sqlstr, u.Subject, u.Name, u.CreatedDate, u.ChangedDate, u.DeletedDate)
	err = db.QueryRow(sqlstr, u.Subject, u.Name, u.CreatedDate, u.ChangedDate, u.DeletedDate).Scan(&u.ID)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// InsertUserByFields inserts the User to the database.
func (s *PostgresStorage) InsertUserByFields(db XODB, u *User) error {
	var err error

	params := make([]interface{}, 0, 5)
	fields := make([]string, 0, 5)
	retCols := `"id"`
	retVars := make([]interface{}, 0, 5)
	retVars = append(retVars, &u.ID)
	fields = append(fields, `"subject"`)
	params = append(params, u.Subject)
	if u.Name.Valid {
		fields = append(fields, `"name"`)
		params = append(params, u.Name)
	} else {
		retCols += `, "name"`
		retVars = append(retVars, &u.Name)
	}
	if u.CreatedDate.Valid {
		fields = append(fields, `"created_date"`)
		params = append(params, u.CreatedDate)
	} else {
		retCols += `, "created_date"`
		retVars = append(retVars, &u.CreatedDate)
	}
	if u.ChangedDate.Valid {
		fields = append(fields, `"changed_date"`)
		params = append(params, u.ChangedDate)
	} else {
		retCols += `, "changed_date"`
		retVars = append(retVars, &u.ChangedDate)
	}
	if u.DeletedDate.Valid {
		fields = append(fields, `"deleted_date"`)
		params = append(params, u.DeletedDate)
	} else {
		retCols += `, "deleted_date"`
		retVars = append(retVars, &u.DeletedDate)
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

	sqlstr := `INSERT INTO "public"."user" (` +
		strings.Join(fields, ",") +
		`) VALUES (` + placeHolders +
		`) RETURNING ` + retCols

	XOLog(sqlstr, params...)
	err = db.QueryRow(sqlstr, params...).Scan(retVars...)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// UpdateUser updates the User in the database.
func (s *PostgresStorage) UpdateUser(db XODB, u *User) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if u._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query

	const sqlstr = `UPDATE "public"."user" SET (` +
		`"subject", "name", "created_date", "changed_date", "deleted_date"` +
		`) = ( ` +
		`$1, $2, $3, $4, $5` +
		`) WHERE "id" = $6`

	// run query
	XOLog(sqlstr, u.Subject, u.Name, u.CreatedDate, u.ChangedDate, u.DeletedDate, u.ID)
	_, err = db.Exec(sqlstr, u.Subject, u.Name, u.CreatedDate, u.ChangedDate, u.DeletedDate, u.ID)
	return err
}

// UpdateUserByFields updates the User in the database.
func (s *PostgresStorage) UpdateUserByFields(db XODB, u *User, fields, retCols []string, params, retVars []interface{}) error {
	var placeHolders string
	for i := range params {
		placeHolders += "$" + strconv.Itoa(i+1)
		if i < len(params)-1 {
			placeHolders += ", "
		}
	}
	params = append(params, u.ID)

	var sqlstr string
	if len(fields) == 1 {
		sqlstr = `UPDATE "public"."user" SET ` +
			strings.Join(fields, ",") +
			` = ` + placeHolders +
			` WHERE id = $` + strconv.Itoa(len(params)) +
			` RETURNING ` + strings.Join(retCols, ", ")
	} else {
		sqlstr = `UPDATE "public"."user" SET (` +
			strings.Join(fields, ",") +
			`) = (` + placeHolders +
			`) WHERE id = $` + strconv.Itoa(len(params)) +
			` RETURNING ` + strings.Join(retCols, ", ")
	}
	XOLog(sqlstr, params...)
	if err := db.QueryRow(sqlstr, params...).Scan(retVars...); err != nil {
		return err
	}

	return nil
}

// SaveUser saves the User to the database.
func (s *PostgresStorage) SaveUser(db XODB, u *User) error {
	if u.Exists() {
		return s.UpdateUser(db, u)
	}

	return s.InsertUser(db, u)
}

// UpsertUser performs an upsert for User.
func (s *PostgresStorage) UpsertUser(db XODB, u *User) error {
	var err error

	// sql query
	const sqlstr = `INSERT INTO "public"."user" (` +
		`"id", "subject", "name", "created_date", "changed_date", "deleted_date"` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) ON CONFLICT ("id") DO UPDATE SET (` +
		`"id", "subject", "name", "created_date", "changed_date", "deleted_date"` +
		`) = (` +
		`EXCLUDED."id", EXCLUDED."subject", EXCLUDED."name", EXCLUDED."created_date", EXCLUDED."changed_date", EXCLUDED."deleted_date"` +
		`)`

	// run query
	XOLog(sqlstr, u.ID, u.Subject, u.Name, u.CreatedDate, u.ChangedDate, u.DeletedDate)
	_, err = db.Exec(sqlstr, u.ID, u.Subject, u.Name, u.CreatedDate, u.ChangedDate, u.DeletedDate)
	if err != nil {
		return err
	}

	// set existence
	u._exists = true

	return nil
}

// DeleteUser deletes the User from the database.
func (s *PostgresStorage) DeleteUser(db XODB, u *User) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return nil
	}

	// if deleted, bail
	if u._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM "public"."user" WHERE "id" = $1`

	// run query
	XOLog(sqlstr, u.ID)
	_, err = db.Exec(sqlstr, u.ID)
	if err != nil {
		return err
	}

	// set deleted
	u._deleted = true

	return nil
}

// DeleteUsers deletes the User from the database.
func (s *PostgresStorage) DeleteUsers(db XODB, us []*User) error {
	var err error

	if len(us) == 0 {
		return nil
	}

	var args []interface{}
	var placeholder string
	for i, u := range us {
		args = append(args, u.ID)
		if i != 0 {
			placeholder = placeholder + ", "
		}
		placeholder += fmt.Sprintf("$%d", i+1)
	}

	// sql query
	var sqlstr = `DELETE FROM "public"."user" WHERE "id" in (` + placeholder + `)`

	// run query
	XOLog(sqlstr, args...)
	_, err = db.Exec(sqlstr, args...)
	if err != nil {
		return err
	}

	// set deleted
	for _, u := range us {
		u._deleted = true
	}

	return nil
}

// GetMostRecentUser returns n most recent rows from 'user',
// ordered by "created_date" in descending order.
func (s *PostgresStorage) GetMostRecentUser(db XODB, n int) ([]*User, error) {
	const sqlstr = `SELECT ` +
		`"id", "subject", "name", "created_date", "changed_date", "deleted_date" ` +
		`FROM "public"."user" ` +
		`ORDER BY created_date DESC LIMIT $1`

	XOLog(sqlstr, n)
	q, err := db.Query(sqlstr, n)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*User
	for q.Next() {
		u := User{}

		// scan
		err = q.Scan(&u.ID, &u.Subject, &u.Name, &u.CreatedDate, &u.ChangedDate, &u.DeletedDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}

// GetMostRecentChangedUser returns n most recent rows from 'user',
// ordered by "changed_date" in descending order.
func (s *PostgresStorage) GetMostRecentChangedUser(db XODB, n int) ([]*User, error) {
	const sqlstr = `SELECT ` +
		`"id", "subject", "name", "created_date", "changed_date", "deleted_date" ` +
		`FROM "public"."user" ` +
		`ORDER BY changed_date DESC LIMIT $1`

	XOLog(sqlstr, n)
	q, err := db.Query(sqlstr, n)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*User
	for q.Next() {
		u := User{}

		// scan
		err = q.Scan(&u.ID, &u.Subject, &u.Name, &u.CreatedDate, &u.ChangedDate, &u.DeletedDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}

// GetAllUser returns all rows from 'user', based on the UserQueryArguments.
// If the UserQueryArguments is nil, it will use the default UserQueryArguments instead.
func (s *PostgresStorage) GetAllUser(db XODB, queryArgs *UserQueryArguments) ([]*User, error) { // nolint: gocyclo
	queryArgs = ApplyUserQueryArgsDefaults(queryArgs)

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
		"name":         true,
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
	params = append(params, *queryArgs.Offset)
	offsetPos := len(params)

	params = append(params, *queryArgs.Limit)
	limitPos := len(params)

	var sqlstr = fmt.Sprintf(`SELECT %s FROM %s WHERE %s deleted_date IS %s ORDER BY %s %s OFFSET $%d LIMIT $%d`,
		`"id", "subject", "name", "created_date", "changed_date", "deleted_date" `,
		`"public"."user"`,
		placeHolders,
		dead,
		orderBy,
		desc,
		offsetPos,
		limitPos)
	XOLog(sqlstr, params...)

	q, err := db.Query(sqlstr, params...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*User
	for q.Next() {
		u := User{}

		// scan
		err = q.Scan(&u.ID, &u.Subject, &u.Name, &u.CreatedDate, &u.ChangedDate, &u.DeletedDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}

// CountAllUser returns a count of all rows from 'user'
func (s *PostgresStorage) CountAllUser(db XODB, queryArgs *UserQueryArguments) (int, error) {
	queryArgs = ApplyUserQueryArgsDefaults(queryArgs)

	dead := "NULL"
	if *queryArgs.Dead {
		dead = "NOT NULL"
	}

	var params []interface{}
	placeHolders := ""

	var err error
	var sqlstr = fmt.Sprintf(`SELECT count(*) from "public"."user" WHERE %s deleted_date IS %s`, placeHolders, dead)
	XOLog(sqlstr)

	var count int
	err = db.QueryRow(sqlstr, params...).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// UsersBySubjectFK retrieves rows from "public"."user" by foreign key Subject.
// Generated from foreign key Account.
func (s *PostgresStorage) UsersBySubjectFK(db XODB, subject string, queryArgs *UserQueryArguments) ([]*User, error) {
	queryArgs = ApplyUserQueryArgsDefaults(queryArgs)

	desc := ""
	if *queryArgs.Desc {
		desc = "DESC"
	}

	dead := "NULL"
	if *queryArgs.Dead {
		dead = "NOT NULL"
	}

	var params []interface{}
	placeHolders := ""
	params = append(params, subject)
	placeHolders = fmt.Sprintf("%s subject = $%d AND ", placeHolders, len(params))

	params = append(params, *queryArgs.Offset)
	offsetPos := len(params)

	params = append(params, *queryArgs.Limit)
	limitPos := len(params)

	var sqlstr = fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s deleted_date IS %s ORDER BY %s %s OFFSET $%d LIMIT $%d`,
		`"id", "subject", "name", "created_date", "changed_date", "deleted_date" `,
		`"public"."user"`,
		placeHolders,
		dead,
		"id",
		desc,
		offsetPos,
		limitPos)

	q, err := db.Query(sqlstr, params...)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*User
	for q.Next() {
		u := User{}

		// scan
		err = q.Scan(&u.ID, &u.Subject, &u.Name, &u.CreatedDate, &u.ChangedDate, &u.DeletedDate)
		if err != nil {
			return nil, err
		}

		res = append(res, &u)
	}

	return res, nil
}

// CountUsersBySubjectFK count rows from "public"."user" by foreign key Subject.
// Generated from foreign key Account.
func (s *PostgresStorage) CountUsersBySubjectFK(db XODB, subject string, queryArgs *UserQueryArguments) (int, error) {
	queryArgs = ApplyUserQueryArgsDefaults(queryArgs)

	dead := "NULL"
	if *queryArgs.Dead {
		dead = "NOT NULL"
	}

	var params []interface{}
	placeHolders := ""
	params = append(params, subject)
	placeHolders = fmt.Sprintf("%s subject = $%d AND ", placeHolders, len(params))

	var err error
	var sqlstr = fmt.Sprintf(`SELECT count(*) from "public"."user" WHERE %s deleted_date IS %s`, placeHolders, dead)
	XOLog(sqlstr)

	var count int
	err = db.QueryRow(sqlstr, params...).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

// AccountInUser returns the Account associated with the User's Subject (subject).
//
// Generated from foreign key 'user_account_subject_fk'.
func (s *PostgresStorage) AccountInUser(db XODB, u *User) (*Account, error) {
	return s.AccountBySubject(db, u.Subject)
}

// UserByID retrieves a row from '"public"."user"' as a User.
//
// Generated from index 'user_pk'.
func (s *PostgresStorage) UserByID(db XODB, id int) (*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`"id", "subject", "name", "created_date", "changed_date", "deleted_date" ` +
		`FROM "public"."user" ` +
		`WHERE "id" = $1`

	// run query
	XOLog(sqlstr, id)
	u := User{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&u.ID, &u.Subject, &u.Name, &u.CreatedDate, &u.ChangedDate, &u.DeletedDate)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
