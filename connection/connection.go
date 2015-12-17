/*
Copyright (c) 2015, Simon J Mudd
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package connection holds information about a connection to the database.
package connection

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	driver = "mysql"
)

// Connection holds the database connection information
type Connection struct {
	dsn     string
	db      *sql.DB
	version string
}

type Connections [](*Connection)

// NewConnection returns a pointer to a new Connection object
func NewConnection(dsn string) *Connection {
	c := new(Connection)

	if c != nil {
		c.dsn = dsn
	}

	return c
}

// Open opens the database using the previously provided dsn.
// function return behaviour is a bit wrong here really. Should pass back to caller.
func (c *Connection) Open() {
	var err error

	c.db, err = sql.Open(driver, c.dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error opening connection to dsn: '%s': %s", c.dsn, err))
	}
	err = c.db.Ping() // force connection to the db.
	if err != nil {
		log.Fatal(fmt.Sprintf("Error pinging host: '%s': %s", c.Hostname(), err))
	}
	c.getVersion()
}

// User() returns the username of the dsn
func (c Connection) User() string {
	var u string

	i := strings.Index(c.dsn, ":")
	if i != -1 {
		u = c.dsn[0:i]
	}

	return u
}

// Hostname() returns the host[:port] combination of the dsn.
// - if the <:port> part is ':3306' then remove it.
// input: user:password@tcp(host.example.com:3306)/performance_schema
// output: host.example.com
func (c Connection) Hostname() string {
	h := c.dsn
	s := strings.Index(h, "(")
	e := strings.Index(h, ")")
	if s != -1 && s != -1 {
		h = h[s+1 : e]
	}
	// remove trailing :3306 from abc:3306 if present
	if len(h) > 5 && h[len(h)-5:] == ":3306" {
		h = h[0 : len(h)-5]
	}
	return h
}

// Handle returns the handle to the database
func (c *Connection) Handle() *sql.DB {
	return c.db
}

func (c Connection) Version() string {
	return c.version
}

// Get the version from MySQL which we may need for later
func (c *Connection) getVersion() {
	var version string

	if c == nil {
		return
	}
	if c.version != "" {
		return // version should never change (unless we're pointing to some sort of pool
	}

	err := c.db.QueryRow("SELECT @@VERSION").Scan(&version)
	switch {
	case err == sql.ErrNoRows:
		// do nothing
	case err != nil:
		log.Fatal(err)
	default:
		c.version = version
	}
}

// return uptime in seconds. Base on collected uptime and collection time.
func (c Connection) Uptime() uint64 {
	var uptime uint64

	err := c.db.QueryRow("SELECT VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME = 'Uptime'").Scan(uptime)
	switch {
	case err == sql.ErrNoRows:
		log.Fatal("No rows when selecting uptime", err)
	case err != nil:
		log.Fatal(err)
	default:
		// do nothing
	}

	return uptime
}

func (c Connection) Name() string {
	return fmt.Sprintf("'%s'@'%s'", c.User(), c.Hostname())
}

// Close closes a single connection
func (c *Connection) Close() {
	if c.db != nil {
		err := c.db.Close()
		if err != nil {
			log.Fatal("Failed to close handle for host:", c.Hostname())
		}
	}
}

// Close closes all valid connections
func (c Connections) Close() {
	if c != nil {
		for i := range c {
			if c[i] != nil {
				c[i].Close()
			}
		}
	}
}
