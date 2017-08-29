package resourceone

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/vincentserpoul/gorestarter/pkg/storage"
)

func TestResourceone_Create(t *testing.T) {
	e := &Resourceone{Label: `test`}
	err := e.Create(context.Background(), pool)
	if err != nil {
		t.Errorf("creation triggered an error %v", err)
		return
	}
	if e.ID == 0 {
		t.Errorf("creation didn't give back an id")
		return
	}
	if e.TimeCreated.Before(time.Now().Add(-1*time.Minute)) ||
		e.TimeUpdated.Before(time.Now().Add(-1*time.Minute)) {
		t.Errorf("creation didn't set the right times for timeCreated and timeUpdated")
		return
	}
}

func BenchmarkResourceone_Create(b *testing.B) {
	for i := 0; i < b.N; i++ {
		e := &Resourceone{Label: `test`}
		err := e.Create(context.Background(), pool)
		if err != nil {
			b.Errorf("Create resourceone error: %v", err)
			return
		}
		// take the opportunity to add some testing IDs
		testResourceoneIDs = append(testResourceoneIDs, e.ID)
	}
}

func TestSelectByID(t *testing.T) {
	ec := &Resourceone{Label: `test`}
	_ = ec.Create(context.Background(), pool)

	e, err := SelectByID(context.Background(), pool, ec.ID)
	if err != nil {
		t.Errorf("SelectByID triggered an error %v", err)
		return
	}
	if e.ID != ec.ID {
		t.Errorf("SelectByID didn't give back the right ID")
		return
	}
	if e.Label != `test` {
		t.Errorf("SelectByID didn't give back the right label")
		return
	}
	if e.TimeCreated.Before(time.Now().Add(-1*time.Minute)) ||
		e.TimeUpdated.Before(time.Now().Add(-1*time.Minute)) {
		t.Errorf("SelectByID didn't get the times")
		return
	}
}

func BenchmarkSelectByID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := SelectByID(context.Background(), pool, testResourceoneIDs[b.N%len(testResourceoneIDs)])
		if err != nil {
			b.Errorf("SelectByID error: %v", err)
			return
		}
	}
}

func TestSelectByTimeUpdated(t *testing.T) {
	ec := &Resourceone{Label: `test`}
	_ = ec.Create(context.Background(), pool)

	es, err := SelectByTimeUpdated(context.Background(), pool, time.Now().Add(-10*time.Second))
	if err != nil {
		t.Errorf("SelectByTimeUpdated triggered an error %v", err)
		return
	}
	if len(es) == 0 {
		t.Errorf("SelectByTimeUpdated didn't give back the right ID")
		return
	}
	for _, e := range es {
		if e.TimeUpdated.Before(time.Now().Add(-8 * time.Second)) {
			t.Errorf("SelectByTimeUpdated didn't get the select the right resourceone entities")
			return
		}
	}
}

func BenchmarkSelectByTimeUpdated(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := SelectByTimeUpdated(context.Background(), pool, time.Now().Add(-1*time.Minute))
		if err != nil {
			b.Errorf("SelectByTimeUpdated error: %v", err)
			return
		}
	}
}

func TestResourceone_Update(t *testing.T) {
	ec := &Resourceone{Label: `test`}
	_ = ec.Create(context.Background(), pool)

	tests := []struct {
		name            string
		id              int64
		wantErr         bool
		wantSQLNotFound bool
		expectedE       *Resourceone
	}{
		{
			name:            "normal existing id",
			id:              ec.ID,
			wantErr:         false,
			wantSQLNotFound: false,
			expectedE: &Resourceone{
				ID:          ec.ID,
				Label:       `testUpdate`,
				TimeCreated: ec.TimeCreated,
				TimeUpdated: time.Now(),
			},
		},
		{
			name:            "not existing existing id",
			id:              999999999,
			wantErr:         true,
			wantSQLNotFound: true,
			expectedE:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Resourceone{Label: `testUpdate`}
			err := Update(context.Background(), pool, tt.id, e)
			if err != nil && !tt.wantErr {
				t.Errorf("Update triggered an error: %v", err)
				return
			}
			if err != nil && tt.wantSQLNotFound && err != ErrSQLNotFound {
				t.Errorf("Update didn't trigger an error where it should")
				return
			}
			if err == nil && tt.wantErr {
				t.Errorf("Update didn't trigger an error where it should")
				return
			}
		})
	}
}

func BenchmarkResourceone_Update(b *testing.B) {
	maxEindex := len(testResourceoneIDs) - 1
	for i := 0; i < b.N; i++ {
		e := &Resourceone{ID: testResourceoneIDs[b.N%maxEindex], Label: `testupdate`}
		err := Update(context.Background(), pool, e.ID, e)
		if err != nil && err != ErrSQLNotFound {
			b.Errorf("Update(%d) error: %v", b.N%maxEindex, err)
			return
		}
	}
}

func TestResourceone_Delete(t *testing.T) {
	ec := &Resourceone{Label: `test`}
	_ = ec.Create(context.Background(), pool)

	err := Delete(context.Background(), pool, ec.ID)
	if err != nil {
		t.Errorf("Delete triggered an error %v", err)
		return
	}

	// Check if DB has been updated as well
	eu, err := SelectByID(context.Background(), pool, ec.ID)
	if err != ErrSQLNotFound && err != nil {
		t.Errorf("Delete: select triggered an error %v", err)
		return
	}
	if eu != nil {
		t.Errorf("Delete didn't delete resourceone in DB")
		return
	}
}

func BenchmarkResourceone_Delete(b *testing.B) {
	maxEindex := len(testResourceoneIDs) - 1
	for i := 0; i < b.N && i <= maxEindex; i++ {
		err := Delete(context.Background(), pool, testResourceoneIDs[i])
		if err != nil && err != ErrSQLNotFound {
			b.Errorf("Delete(%d) error: %v", b.N%maxEindex, err)
			return
		}
	}
}

var pool *sqlx.DB
var testResourceoneIDs []int64

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	newConnPool, err := storage.NewMySQLDBConnPool(&storage.MySQLDBConf{
		Protocol: "tcp",
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "internal",
		Password: "dev",
		DbName:   "test",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		errClose := newConnPool.Close()
		if errClose != nil {
			log.Fatalf("%v", errClose)
		}
	}()

	pool = newConnPool

	ddls := []DDL{DDL{}}

	for _, ddl := range ddls {
		errD := ddl.MigrateDown(ctx, pool)
		if errD != nil {
			log.Fatal(errD)
		}
		errU := ddl.MigrateUp(ctx, pool)
		if errU != nil {
			log.Fatal(errU)
		}
	}

	retCode := m.Run()

	for _, ddl := range ddls {
		errD := ddl.MigrateDown(ctx, pool)
		if errD != nil {
			log.Fatal(errD)
		}
	}

	os.Exit(retCode)

}
