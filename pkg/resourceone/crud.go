package resourceone

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// ErrSQLNotFound is returned when no rows affected or found
var ErrSQLNotFound = errors.New("no resourceone found")

// Create will create an resourceone in the DB
func (e *Resourceone) Create(
	ctx context.Context,
	db *sqlx.DB,
) error {
	// Insert in DB
	resourceoneID, errIns := insertOne(ctx, db, e.Label)
	if errIns != nil {
		return fmt.Errorf("Create(%s): %v", e.Label, errIns)
	}

	e.ID = resourceoneID
	e.TimeCreated = time.Now()
	e.TimeUpdated = time.Now()

	return nil
}

func insertOne(
	ctx context.Context,
	db *sqlx.DB,
	label string,
) (int64, error) {

	res, err := db.NamedExecContext(
		ctx,
		`
			INSERT INTO resourceone(label)
			VALUES (:label)
		`,
		map[string]interface{}{
			"label": label,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("insertOne(%s): %v", label, err)
	}

	id, errL := res.LastInsertId()
	if errL != nil {
		return 0, fmt.Errorf("insertOne(%s): %v", label, errL)
	}

	return id, nil
}

// queryFilter allows for select filters
type queryFilter struct {
	filterSQL   string
	namedParams map[string]interface{}
}

// filterByID will filter the select query by resourceone.ID
func filterByID(resourceoneID int64) queryFilter {
	return queryFilter{
		filterSQL: " AND resourceone_id = :resourceoneID ",
		namedParams: map[string]interface{}{
			"resourceoneID": resourceoneID,
		},
	}
}

// SelectByID returns one resourceone entity
func SelectByID(
	ctx context.Context,
	db *sqlx.DB,
	resourceoneID int64,
) (*Resourceone, error) {
	es, err := selectsql(ctx, db, filterByID(resourceoneID))
	if err != nil {
		return nil, fmt.Errorf("SelectByID(%d): %v", resourceoneID, err)
	}

	if len(es) == 0 {
		return nil, ErrSQLNotFound
	}

	return es[0], nil
}

// FilterByTimeUpdated will filter the select query by resourceone.ID
func filterByTimeUpdated(updatedAfter time.Time) queryFilter {
	return queryFilter{
		filterSQL: " AND time_updated > :updatedAfter ",
		namedParams: map[string]interface{}{
			"updatedAfter": updatedAfter,
		},
	}
}

// SelectByTimeUpdated will get all the entityone updated after a certain date
func SelectByTimeUpdated(
	ctx context.Context,
	db *sqlx.DB,
	updatedAfter time.Time,
) ([]*Resourceone, error) {
	es, err := selectsql(ctx, db, filterByTimeUpdated(updatedAfter))
	if err != nil {
		return nil, fmt.Errorf("SelectByTimeUpdated(%v): %v", updatedAfter, err)
	}

	if len(es) == 0 {
		return nil, ErrSQLNotFound
	}

	return es, nil
}

// selectsql will get a specific resourceone from the DB and return an resourceone struct
func selectsql(
	ctx context.Context,
	db *sqlx.DB,
	queryFilters ...queryFilter,
) ([]*Resourceone, error) {

	query := `SELECT resourceone_id, label, time_created, time_updated
				FROM resourceone
				WHERE 0=0 `
	namedParams := make(map[string]interface{})

	// merge filters into the query
	for _, filter := range queryFilters {
		query += filter.filterSQL
		for k, v := range filter.namedParams {
			if _, ok := filter.namedParams[k]; ok {
				namedParams[k] = v
			}
		}
	}

	rows, err := db.NamedQueryContext(ctx, query, namedParams)
	if err != nil {
		return nil, fmt.Errorf("Select(%v): %v", queryFilters, err)
	}

	var es []*Resourceone

	for rows.Next() {
		e := &Resourceone{}
		err := rows.StructScan(e)
		if err != nil {
			return nil, fmt.Errorf("Select(%v): %v", queryFilters, err)
		}
		es = append(es, e)
	}

	return es, nil
}

// Update will update an specific resourceone in the DB
func Update(
	ctx context.Context,
	db *sqlx.DB,
	resourceoneID int64,
	e *Resourceone,
) error {

	res, err := db.NamedExecContext(
		ctx,
		`
			UPDATE resourceone
				SET label = :label,
					time_updated = NOW()
			WHERE resourceone_id = :resourceoneID
		`,
		map[string]interface{}{
			"label":         e.Label,
			"resourceoneID": resourceoneID,
		},
	)
	if err != nil {
		return fmt.Errorf("Update(%d, %s): %v", e.ID, e.Label, err)
	}

	ra, errRA := res.RowsAffected()
	if errRA != nil {
		return fmt.Errorf("Delete(%d): %v", resourceoneID, errRA)
	}

	if ra == 0 {
		return ErrSQLNotFound
	}

	return nil
}

// Delete will delete an resourceone from the DB
func Delete(
	ctx context.Context,
	db *sqlx.DB,
	resourceoneID int64,
) error {

	res, err := db.NamedExecContext(
		ctx,
		`
			DELETE FROM resourceone
			WHERE resourceone_id = :resourceoneID
		`,
		map[string]interface{}{
			"resourceoneID": resourceoneID,
		},
	)
	if err != nil {
		return fmt.Errorf("Delete(%d): %v", resourceoneID, err)
	}

	ra, errRA := res.RowsAffected()
	if errRA != nil {
		return fmt.Errorf("Delete(%d): %v", resourceoneID, errRA)
	}

	if ra == 0 {
		return ErrSQLNotFound
	}

	return nil
}
