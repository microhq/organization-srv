package handler

import (
	"strings"

	"github.com/micro/go-micro/errors"
	"github.com/microhq/organization-srv/db"
	"github.com/microhq/organization-srv/proto/org"
	"golang.org/x/net/context"
)

func (s *Org) CreateMember(ctx context.Context, req *org.CreateMemberRequest, rsp *org.CreateMemberResponse) error {
	req.Member.Username = strings.ToLower(req.Member.Username)
	if err := db.CreateMember(req.Member); err != nil {
		return errors.InternalServerError("go.micro.srv.organization.CreateMember", err.Error())
	}
	return nil
}

func (s *Org) ReadMember(ctx context.Context, req *org.ReadMemberRequest, rsp *org.ReadMemberResponse) error {
	member, err := db.ReadMember(req.Id)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.organization.ReadMember", err.Error())
	}
	rsp.Member = member
	return nil
}

func (s *Org) UpdateMember(ctx context.Context, req *org.UpdateMemberRequest, rsp *org.UpdateMemberResponse) error {
	req.Member.Username = strings.ToLower(req.Member.Username)
	if err := db.UpdateMember(req.Member); err != nil {
		return errors.InternalServerError("go.micro.srv.organization.UpdateMember", err.Error())
	}
	return nil
}

func (s *Org) DeleteMember(ctx context.Context, req *org.DeleteMemberRequest, rsp *org.DeleteMemberResponse) error {
	if err := db.DeleteMember(req.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.organization.DeleteMember", err.Error())
	}
	return nil
}

func (s *Org) SearchMembers(ctx context.Context, req *org.SearchMembersRequest, rsp *org.SearchMembersResponse) error {
	members, err := db.SearchMembers(req.OrgName, req.Username, req.Limit, req.Offset)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.organization.SearchMember", err.Error())
	}
	rsp.Members = members
	return nil
}
