package main

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.SetDialect("mysql")
	goose.AddMigration(upUsersTable, downUsersTable)
}

func upUsersTable(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	sql := `
        CREATE TABLE IF NOT EXISTS users (
            id int(11) not null auto_increment,
            username varchar(100) not null,
            primary key(id)
        )ENGINE=InnoDB
    `
	_, err := tx.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func downUsersTable(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	sql := `
        DROP TABLE IF EXISTS users
    `
	_, err := tx.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
