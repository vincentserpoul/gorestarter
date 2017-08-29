package storage

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestNewMySQLDBConnPool(t *testing.T) {

	tests := []struct {
		name        string
		mysqlDBConf *MySQLDBConf
		wantErr     bool
	}{
		{
			name: "Working connection",
			mysqlDBConf: &MySQLDBConf{
				Protocol: "tcp",
				Host:     "127.0.0.1",
				Port:     "3306",
				User:     "internal",
				Password: "dev",
				DbName:   "test",
			},
			wantErr: false,
		},
		{
			name: "Working connection",
			mysqlDBConf: &MySQLDBConf{
				Protocol: "tcp",
				Host:     "none",
				Port:     "3306",
				User:     "internal",
				Password: "badtest",
				DbName:   "badtest",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMySQLDBConnPool(tt.mysqlDBConf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMySQLDBConnPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
