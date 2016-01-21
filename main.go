package main

import (
	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/organization-srv/db"
	"github.com/micro/organization-srv/handler"
	"github.com/micro/organization-srv/proto/org"
)

func main() {

	service := micro.NewService(
		micro.Name("go.micro.srv.organization"),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/organization",
			},
		),
		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				db.Url = c.String("database_url")
			}
		}),
	)

	service.Init()

	db.Init()

	org.RegisterOrgHandler(service.Server(), new(handler.Org))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
