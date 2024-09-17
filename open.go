package lqs

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type connector struct {
	driver     driver.Driver
	dsn        string
	preemptive string
}

var _ driver.Connector = &connector{}

// Connect implements the logic of returns a new connection
// with the preemptive SQL code run.
func (c *connector) Connect(context.Context) (driver.Conn, error) {
	conn, err := c.driver.Open(c.dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open dsn: %w", err)
	}

	stmt, err := conn.Prepare(c.preemptive)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("could not set preemptive: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(nil)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("could not execute preemptive: %w", err)
	}

	return conn, nil
}

// Driver just returns the driver.
func (c *connector) Driver() driver.Driver {
	return c.driver
}

// Open takes the same parameters as `sql.Open` and returns the same values.
// The extra parameter `preemptive` is SQL statements to be run when a new client is created.
func Open(driverName string, dsn string, preemptive string) (*sql.DB, error) {
	db, err := sql.Open(driverName, "")
	if err != nil {
		return nil, fmt.Errorf("could not load db: %w", err)
	}
	defer db.Close()

	connector := &connector{
		driver:     db.Driver(),
		dsn:        dsn,
		preemptive: preemptive,
	}

	return sql.OpenDB(connector), nil
}
