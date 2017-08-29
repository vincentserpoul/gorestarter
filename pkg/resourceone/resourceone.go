package resourceone

import "time"

// Resourceone represents an entity
type Resourceone struct {
	ID          int64     `db:"resourceone_id" json:"resourceoneId"`
	Label       string    `db:"label" json:"label"`
	TimeCreated time.Time `db:"time_created" json:"timeCreated"`
	TimeUpdated time.Time `db:"time_updated" json:"timeUpdated"`
}
