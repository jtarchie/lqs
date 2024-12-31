package lqs_test

import (
	"context"
	"os"
	"runtime"
	"sync"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/jtarchie/lqs"
	_ "github.com/mattn/go-sqlite3"
)

func TestPragma(t *testing.T) {
	tempFile, err := os.CreateTemp("", "lqs_test")
	assert.NoError(t, err)

	client, err := lqs.Open("sqlite3", "file://"+tempFile.Name(), "PRAGMA cache_size = 1234;")
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

func TestThreadSafe(t *testing.T) {
	tempFile, err := os.CreateTemp("", "lqs_test")
	assert.NoError(t, err)

	client, err := lqs.Open("sqlite3", "file://"+tempFile.Name(), "PRAGMA cache_size = 1234;")
	assert.NoError(t, err)
	defer client.Close()

	client.SetMaxOpenConns(runtime.NumCPU())
	client.SetMaxIdleConns(runtime.NumCPU())

	var wg sync.WaitGroup

	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var cache_size int64

			err := client.QueryRow(`PRAGMA cache_size;`).Scan(&cache_size)
			assert.NoError(t, err)
			assert.Equal(t, cache_size, 1234)

			rows, err := client.Query("PRAGMA cache_size;")
			assert.NoError(t, err)
			defer rows.Close()

			cache_size = 0
			rows.Next()
			err = rows.Scan(&cache_size)
			assert.NoError(t, err)
			assert.Equal(t, cache_size, 1234)
		}()
	}
	wg.Wait()
}

func TestForeignKeys(t *testing.T) {
	client, err := lqs.Open("sqlite3", ":memory:", "PRAGMA foreign_keys = ON;")
	assert.NoError(t, err)
	defer client.Close()

	var foreign_keys int64
	err = client.QueryRow(`PRAGMA foreign_keys;`).Scan(&foreign_keys)
	assert.NoError(t, err)
	assert.Equal(t, foreign_keys, 1)
}

// ensure errors are propagated correctly
func TestUnknownDriver(t *testing.T) {
	_, err := lqs.Open("blahblah", "", "")
	assert.Error(t, err)
}
