package store

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
)

type Model interface {
	Get(client *Client) error
	Delete(client *Client) error
	Update(client *Client) error
	Create(client *Client) error
}

type Query interface {
	Get(string) string
}

type Collection interface {
	Get(client *Client, query Query) error
}

type Client struct {
	Ctxt         context.Context
	port         int
	hostname     string
	username     string
	password     string
	databaseName string
	tLS          bool
	db           *pgx.Conn
}

func New(ctxt context.Context, userName, hostName, password, dataBaseName string, port int, tls bool) *Client {
	return &Client{ctxt, port, hostName, userName, password, dataBaseName, tls, nil}
}

func (client *Client) Open() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", client.username, client.password, client.hostname, client.port, client.databaseName)
	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		return err
	}

	client.db = conn

	return nil
}

func (client *Client) Close() error {
	return client.db.Close(context.Background())
}
