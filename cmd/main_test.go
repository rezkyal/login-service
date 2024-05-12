package main

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_newServer(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	os.Setenv("DATABASE_URL", "sqlmock_db_0")

	got := newServer()

	assert.NotNil(t, got)
}
