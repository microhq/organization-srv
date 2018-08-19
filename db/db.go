package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/microhq/organization-srv/proto/org"
)

var (
	db                 *sql.DB
	Url                = "root@tcp(127.0.0.1:3306)/organization"
	database           string
	organizationSchema = `CREATE TABLE IF NOT EXISTS organizations (
id varchar(36) primary key,
name varchar(255),
email varchar(255),
owner varchar(255),
created integer,
updated integer,
unique (name));`
	membersSchema = `CREATE TABLE IF NOT EXISTS organization_members (
id varchar(36) primary key,
organization_name varchar(36) ,
username varchar(36),
roles text,
created integer,
updated integer,
unique (organization_name, username));`

	q = map[string]string{
		"delete": "DELETE from %s.%s where id = ?",
		"create": `INSERT into %s.%s (
				id, name, email, owner, created, updated) 
				values (?, ?, ?, ?, ?, ?)`,
		"update":             "UPDATE %s.%s set name = ?, email = ?, owner = ?, updated = ? where id = ?",
		"read":               "SELECT * from %s.%s where id = ?",
		"list":               "SELECT * from %s.%s limit ? offset ?",
		"searchName":         "SELECT * from %s.%s where name = ? limit ? offset ?",
		"searchOwner":        "SELECT * from %s.%s where owner = ? limit ? offset ?",
		"searchNameAndOwner": "SELECT * from %s.%s where name = ? and owner = ? limit ? offset ?",
	}

	mq = map[string]string{
		"deleteMember": "DELETE from %s.%s where id = ?",
		"createMember": `INSERT into %s.%s (
				id, organization_name, username, roles, created, updated) 
				values (?, ?, ?, ?, ?, ?)`,
		"updateMember":         "UPDATE %s.%s set roles = ?, updated = ? where id = ?",
		"readMember":           "SELECT * from %s.%s where id = ?",
		"searchUsername":       "SELECT * from %s.%s where username = ? limit ? offset ?",
		"searchOrg":            "SELECT * from %s.%s where organization_name = ? limit ? offset ?",
		"searchUsernameAndOrg": "SELECT * from %s.%s where organization_name = ? and username = ? limit ? offset ?",
	}

	st = map[string]*sql.Stmt{}
)

func Init() {
	var d *sql.DB
	var err error

	parts := strings.Split(Url, "/")
	if len(parts) != 2 {
		panic("Invalid database url")
	}

	if len(parts[1]) == 0 {
		panic("Invalid database name")
	}

	url := parts[0]
	database = parts[1]

	if d, err = sql.Open("mysql", url+"/"); err != nil {
		log.Fatal(err)
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		log.Fatal(err)
	}
	d.Close()
	if d, err = sql.Open("mysql", Url); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(organizationSchema); err != nil {
		log.Fatal(err)
	}
	if _, err = d.Exec(membersSchema); err != nil {
		log.Fatal(err)
	}
	db = d

	for query, statement := range q {
		prepared, err := db.Prepare(fmt.Sprintf(statement, database, "organizations"))
		if err != nil {
			log.Fatal(err)
		}
		st[query] = prepared
	}
	for query, statement := range mq {
		prepared, err := db.Prepare(fmt.Sprintf(statement, database, "organization_members"))
		if err != nil {
			log.Fatal(err)
		}
		st[query] = prepared
	}
}

func Create(organization *org.Organization) error {
	organization.Created = time.Now().Unix()
	organization.Updated = time.Now().Unix()
	_, err := st["create"].Exec(organization.Id, organization.Name, organization.Email, organization.Owner, organization.Created, organization.Updated)
	return err
}

func Delete(id string) error {
	_, err := st["delete"].Exec(id)
	return err
}

func Update(organization *org.Organization) error {
	organization.Updated = time.Now().Unix()
	_, err := st["update"].Exec(organization.Name, organization.Email, organization.Owner, organization.Updated, organization.Id)
	return err
}

func Read(id string) (*org.Organization, error) {
	organization := &org.Organization{}

	r := st["read"].QueryRow(id)
	if err := r.Scan(&organization.Id, &organization.Name, &organization.Email, &organization.Owner, &organization.Created, &organization.Updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return organization, nil
}

func Search(name, owner string, limit, offset int64) ([]*org.Organization, error) {
	var r *sql.Rows
	var err error

	if len(name) > 0 && len(owner) > 0 {
		r, err = st["searchNameAndOwner"].Query(name, owner, limit, offset)
	} else if len(name) > 0 {
		r, err = st["searchName"].Query(name, limit, offset)
	} else if len(owner) > 0 {
		r, err = st["searchOwner"].Query(owner, limit, offset)
	} else {
		r, err = st["list"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var organizations []*org.Organization

	for r.Next() {
		organization := &org.Organization{}
		if err := r.Scan(&organization.Id, &organization.Name, &organization.Email, &organization.Owner, &organization.Created, &organization.Updated); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		organizations = append(organizations, organization)

	}
	if r.Err() != nil {
		return nil, err
	}

	return organizations, nil
}
