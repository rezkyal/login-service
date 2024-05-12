// This file contains the repository implementation layer.
package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	type args struct {
		opts NewRepositoryOptions
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				opts: NewRepositoryOptions{
					Dsn: "sqlmock_db_0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRepository(tt.args.opts)

			assert.NotNil(t, got)
		})
	}
}
