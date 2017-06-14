// Package zero provides a example schema for a demo
//
// Source: https://github.com/graphql/graphql.github.io/blob/source/site/_core/swapiSchema.js
package zero

import (
	"database/sql"
	"fmt"

	graphql "github.com/neelance/graphql-go"
)

//Schema is the graphql schema
var Schema = `
	schema {
		query: Query
	}
	# The query type, represents all of the entry points into our object graph
	type Query {
		person(id: ID!): Person
	}
	type Person {
		id: ID!
		first_name: String!
		last_name: String!
		username: String!
		email: String!
		friends: [Person]
	}
`

type person struct {
	ID        graphql.ID
	FirstName string
	LastName  string
	Username  string
	Email     string
	Friends   *[]person
}

var personData = make(map[graphql.ID]*person)

//Resolver maps the schema to go
type Resolver struct {
	DB *sql.DB
}

//Person represents the Person type
func (r *Resolver) Person(args struct{ ID graphql.ID }) *personResolver {

	var p person
	err := r.DB.QueryRow(`
		SELECT to_char(uid, '999'), first_name, last_name,
			username, email FROM userinfo where username=$1`,
		args.ID).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Username, &p.Email)
	CheckErr(err)
	return &personResolver{&p, r.DB}

}

type personResolver struct {
	p  *person
	db *sql.DB
}

func (r *personResolver) ID() graphql.ID {
	return r.p.ID
}

func (r *personResolver) FirstName() string {
	return r.p.FirstName
}

func (r *personResolver) LastName() string {
	return r.p.LastName
}

func (r *personResolver) Username() string {
	return r.p.Username
}

func (r *personResolver) Email() string {
	return r.p.Email
}

func (r *personResolver) Friends() *[]*personResolver {
	return resolvePersons(r.db, r.p.Friends)
}

func resolvePersons(db *sql.DB, friends *[]person) *[]*personResolver {
	var persons []*personResolver
	rows, err := db.Query(
		`SELECT to_char(uid, '999'), first_name, last_name,
		username, email FROM userinfo`)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var p person
		err = rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Username, &p.Email)
		if err != nil {
			fmt.Println(err)
		}
		persons = append(persons, &personResolver{&p, db})
	}
	return &persons
}
