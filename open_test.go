package lqs_test

import (
	"context"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/jtarchie/lqs"
	_ "github.com/mattn/go-sqlite3"
)

func TestPragma(t *testing.T) {
	client, err := lqs.Open("sqlite3", ":memory:", "PRAGMA cache_size = 1234;")
	assert.NoError(t, err)
	defer client.Close()

	var cache_size int64
	err = client.QueryRow(`PRAGMA cache_size;`).Scan(&cache_size)
	assert.NoError(t, err)
	assert.Equal(t, cache_size, 1234)

	conn, err := client.Conn(context.Background())
	assert.NoError(t, err)
	rows, err := conn.QueryContext(context.Background(), "PRAGMA cache_size;")
	assert.NoError(t, err)
	defer rows.Close()

	cache_size = 0
	rows.Next()
	err = rows.Scan(&cache_size)
	assert.NoError(t, err)
	assert.Equal(t, cache_size, 1234)
}

// ensure errors are propagated correctly
func TestUnknownDriver(t *testing.T) {
	_, err := lqs.Open("blahblah", "", "")
	assert.Error(t, err)
}
