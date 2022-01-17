package store

import (
	"sync"

	"github.com/jackc/pgx"
)

type Client struct {
	ConnConfig pgx.ConnConfig
	Conn       *pgx.Conn
}

var client *Client

func New(userName, hostName, password, dataBaseName string, port int) *Client {
	var once sync.Once

	once.Do(
		func() {
			client = &Client{
				ConnConfig: pgx.ConnConfig{
					Host: hostName, Port: uint16(port), User: userName, Password: password, Database: dataBaseName,
				}, Conn: nil,
			}
		},
	)
	return client
}

func (client *Client) Open() error {
	conn, err := pgx.Connect(client.ConnConfig)
	if err != nil {
		return err
	}

	client.Conn = conn

	return nil
}

func (client *Client) Close() error {
	return client.Conn.Close()
}
