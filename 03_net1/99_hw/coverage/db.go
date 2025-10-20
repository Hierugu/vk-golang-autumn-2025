package main

import (
	"encoding/xml"
	"os"
)

type XMLUser struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

func (u XMLUser) ToUser() User {
	return User{u.ID, u.FirstName + u.LastName, u.Age, u.About, u.Gender}
}

type Database interface {
	Load() ([]User, error)
}

type XMLDatabase struct {
	filePath string
}

func (db XMLDatabase) Load() ([]User, error) {
	xmlData, err := os.Open(db.filePath)
	if err != nil {
		return nil, err
	}
	defer xmlData.Close()

	rawUser := XMLUser{}
	users := []User{}
	d := xml.NewDecoder(xmlData)
	for t, _ := d.Token(); t != nil; t, _ = d.Token() {
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "row" {
				d.DecodeElement(&rawUser, &se)
				users = append(users, rawUser.ToUser())
			}
		}
	}
	return users, nil
}
