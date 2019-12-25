// Package storage contains the types for schema 'dbo'.
package storage

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
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
	XOLog(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate)
	err = db.QueryRow(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate).Scan(&a.ID)
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
	XOLog(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate, a.ID)
	_, err = db.Exec(sqlstr, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate, a.ID)
	return err
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
	XOLog(sqlstr, a.ID, a.Subject, a.Email, a.CreatedDate, a.ChangedDate, a.DeletedDate)
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
	XOLog(sqlstr, a.ID)
	_, err = db.Exec(sqlstr, a.ID)
	if err != nil {
		return err
	}

	// set deleted
	a._deleted = true

	return nil
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
	XOLog(sqlstr, id)
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
	XOLog(sqlstr, subject)
	a := Account{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, subject).Scan(&a.ID, &a.Subject, &a.Email, &a.CreatedDate, &a.ChangedDate, &a.DeletedDate)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
