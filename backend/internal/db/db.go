package db

import (
	"app/ent"
	"context"
	"log"

	_ "github.com/lib/pq"
)

type Db struct {
	Client *ent.Client
}

func NewDbClient(urlDb string) *Db {
	client, err := ent.Open("postgres", urlDb)
	if err != nil {
		log.Fatal("Error to open connect with postgres", err.Error())
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatal("Error to create migrations", err.Error())
	}

	return &Db{Client: client}
}

func (d *Db) Close() {
	d.Client.Close()
}
