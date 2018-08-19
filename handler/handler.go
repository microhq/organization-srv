package handler

import (
	"strings"

	"github.com/micro/go-micro/errors"
	"github.com/microhq/organization-srv/db"
	"github.com/microhq/organization-srv/proto/org"
	"golang.org/x/net/context"
)

type Org struct{}

func (s *Org) Create(ctx context.Context, req *org.CreateRequest, rsp *org.CreateResponse) error {
	req.Organization.Name = strings.ToLower(req.Organization.Name)
	req.Organization.Email = strings.ToLower(req.Organization.Email)
	req.Organization.Owner = strings.ToLower(req.Organization.Owner)
	if err := db.Create(req.Organization); err != nil {
		return errors.InternalServerError("go.micro.srv.organization.Create", err.Error())
	}
	return nil
}

func (s *Org) Read(ctx context.Context, req *org.ReadRequest, rsp *org.ReadResponse) error {
	org, err := db.Read(req.Id)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.organization.Read", err.Error())
	}
	rsp.Organization = org
	return nil
}

func (s *Org) Update(ctx context.Context, req *org.UpdateRequest, rsp *org.UpdateResponse) error {
	req.Organization.Name = strings.ToLower(req.Organization.Name)
	req.Organization.Email = strings.ToLower(req.Organization.Email)
	req.Organization.Owner = strings.ToLower(req.Organization.Owner)
	if err := db.Update(req.Organization); err != nil {
		return errors.InternalServerError("go.micro.srv.organization.Update", err.Error())
	}
	return nil
}

func (s *Org) Delete(ctx context.Context, req *org.DeleteRequest, rsp *org.DeleteResponse) error {
	if err := db.Delete(req.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.organization.Delete", err.Error())
	}
	return nil
}

func (s *Org) Search(ctx context.Context, req *org.SearchRequest, rsp *org.SearchResponse) error {
	orgs, err := db.Search(req.Name, req.Owner, req.Limit, req.Offset)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.organization.Search", err.Error())
	}
	rsp.Organizations = orgs
	return nil
}
