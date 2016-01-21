package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/micro/organization-srv/proto/org"
)

func CreateMember(member *org.Member) error {
	member.Created = time.Now().Unix()
	member.Updated = time.Now().Unix()
	b, err := json.Marshal(member.Roles)
	if err != nil {
		return err
	}
	_, err = st["createMember"].Exec(member.Id, member.OrgName, member.Username, string(b), member.Created, member.Updated)
	return err
}

func DeleteMember(id string) error {
	_, err := st["deleteMember"].Exec(id)
	return err
}

func UpdateMember(member *org.Member) error {
	b, err := json.Marshal(member.Roles)
	if err != nil {
		return err
	}
	_, err = st["updateMember"].Exec(string(b), time.Now().Unix(), member.Id)
	return err
}

func ReadMember(id string) (*org.Member, error) {
	member := &org.Member{}

	r := st["readMember"].QueryRow(id)
	var roles string
	if err := r.Scan(&member.Id, &member.OrgName, &member.Username, &roles, &member.Created, &member.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	if err := json.Unmarshal([]byte(roles), &member.Roles); err != nil {
		return nil, err
	}

	return member, nil
}

func SearchMembers(orgName, username string, limit, offset int64) ([]*org.Member, error) {
	var r *sql.Rows
	var err error

	if len(orgName) > 0 && len(username) > 0 {
		r, err = st["searchUsernameAndOrg"].Query(orgName, username, limit, offset)
	} else if len(username) > 0 {
		r, err = st["searchUsername"].Query(username, limit, offset)
	} else if len(orgName) > 0 {
		r, err = st["searchOrg"].Query(orgName, limit, offset)
	} else {
		return nil, errors.New("org id and username cannot be blank")
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var members []*org.Member

	for r.Next() {
		member := &org.Member{}
		var roles string
		if err := r.Scan(&member.Id, &member.OrgName, &member.Username, &roles, &member.Created, &member.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		if err := json.Unmarshal([]byte(roles), &member.Roles); err != nil {
			return nil, err
		}
		members = append(members, member)

	}
	if r.Err() != nil {
		return nil, err
	}

	return members, nil
}
