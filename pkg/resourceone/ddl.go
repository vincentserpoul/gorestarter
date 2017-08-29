package resourceone

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// DDL is used to do modifications in the DB
type DDL struct{}

// MigrateUp creates the needed tables
func (ddl *DDL) MigrateUp(ctx context.Context, db *sqlx.DB) error {
	_, errExec := db.ExecContext(
		ctx,
		`
            CREATE TABLE IF NOT EXISTS resourceone (
                resourceone_id BIGINT NOT NULL AUTO_INCREMENT,
                label VARCHAR(50),
                time_created DATETIME NOT NULL DEFAULT NOW(),
				time_updated DATETIME NOT NULL DEFAULT NOW(),
                PRIMARY KEY (resourceone_id),
				INDEX r_tu_idx (time_updated ASC)
            );
    `)

	return errExec
}

// MigrateDown destroys the needed tables
func (ddl *DDL) MigrateDown(ctx context.Context, db *sqlx.DB) error {
	_, errExec := db.ExecContext(
		ctx,
		`
        DROP TABLE IF EXISTS resourceone;
    `)

	return errExec
}
