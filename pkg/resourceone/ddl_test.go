package resourceone

import (
	"context"
	"testing"
)

func TestDDL_MigrateDown(t *testing.T) {
	tests := []struct {
		name            string
		withEmptySchema bool
		wantErr         bool
	}{
		{
			name:            "Default",
			withEmptySchema: true,
			wantErr:         false,
		},
		{
			name:            "Default",
			withEmptySchema: false,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		ddl := &DDL{}
		t.Run(tt.name, func(t *testing.T) {
			if !tt.withEmptySchema {
				errE := ddl.MigrateUp(context.Background(), pool)
				if errE != nil {
					t.Errorf("DDL.MigrateDown() error = %v when emptying the schema", errE)
				}
			}
			if err := ddl.MigrateDown(context.Background(), pool); (err != nil) != tt.wantErr {
				t.Errorf("DDL.MigrateUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// This test needs to be run last, so that the benchmark still have the right tables to run
func TestDDL_MigrateUp(t *testing.T) {

	tests := []struct {
		name            string
		withEmptySchema bool
		wantErr         bool
	}{
		{
			name:            "Default",
			withEmptySchema: true,
			wantErr:         false,
		},
		{
			name:            "Default",
			withEmptySchema: false,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		ddl := &DDL{}
		t.Run(tt.name, func(t *testing.T) {
			if tt.withEmptySchema {
				errE := ddl.MigrateDown(context.Background(), pool)
				if errE != nil {
					t.Errorf("DDL.MigrateDown() error = %v when emptying the schema", errE)
				}
			}
			if err := ddl.MigrateUp(context.Background(), pool); (err != nil) != tt.wantErr {
				t.Errorf("DDL.MigrateUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
