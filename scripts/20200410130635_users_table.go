package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upUsersTable, downUsersTable)
}

func upUsersTable(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	sql := `
        CREATE TABLE users IF NOT EXISTS (
            id int(11) not null auto_increment,
            email varchar(100) not null,
            password varchar(100) not null,
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
