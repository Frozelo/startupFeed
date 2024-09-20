package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

type Postgres struct {
	Conn         *pgx.Conn
	connAttempts int
	connTimeout  time.Duration
}

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = time.Second
)

func New(connString string) (*Postgres, error) {
	pg := &Postgres{
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	var err error
	for pg.connAttempts > 0 {
		pg.Conn, err = pgx.Connect(context.Background(), connString)
		if err == nil {
			break
		}
		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts-1)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - connAttempts == 0: %w", err)
	}
	
	return pg, nil
}

func (s *Postgres) Close() {
	if err := s.Conn.Close(context.Background()); err != nil {
		log.Printf("Postgres - Close: %s", err)
	}
}
