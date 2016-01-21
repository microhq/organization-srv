# Organization Server

Organization server is a service to track orgs and members.

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a mysql database

4. Download and start the service

	```shell
	go get github.com/micro/organization-srv
	organization-srv --database_url="root:root@tcp(192.168.99.100:3306)/organization"
	```

	OR as a docker container

	```shell
	docker run microhq/organization-srv --database_url="root:root@tcp(192.168.99.100:3306)/organization" --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Organization server implements the following RPC Methods

Org
- Create
- Read
- Update
- Delete
- Search
- CreateMember
- ReadMember
- UpdateMember
- DeleteMember
- SearchMembers


### Org.Create
```shell
micro query go.micro.srv.organization Org.Create '{ "organization": {"id": "1", "name": "Micro", "email": "micro@example.com", "owner": "asim"}}'
{}
```

### Org.Search
```shell
micro query go.micro.srv.organization Org.Search '{"limit": 10}'
{
	"organizations": [
		{
			"created": 1.453407783e+09,
			"email": "micro@example.com",
			"id": "1",
			"name": "micro",
			"owner": "asim",
			"updated": 1.453407783e+09
		}
	]
}
```
