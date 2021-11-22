package models

import (
	"golang-gogin/forms"
	"golang-gogin/helpers"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID         bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string        `json:"name" bson:"name"`
	Email      string        `json:"email" bson:"email"`
	Password   string        `json:"password" bson:"password"`
	IsVerified bool          `json:"is_verified" bson:"is_verified"`
}

type UserModel struct{}

func (u *UserModel) Signup(data forms.SignupUsercommand) error {
	collection := dbConnect.Use(databaseName, "user")
	err := collection.Insert(bson.M{
		"name":        data.Name,
		"email":       data.Email,
		"password":    helpers.GeneratePasswordHash([]byte(data.Password)),
		"is_verified": false,
	})

	return err
}

func (u *UserModel) GetUserByEmail(email string) (user User, err error) {
	collection := dbConnect.Use(databaseName, "user")

	err = collection.Find(bson.M{"email": email}).One(&user)

	return user, err
}

func (u *UserModel) UpdateUserPass(email string, hashedPassword string) error {
	collection := dbConnect.Use(databaseName, "user")

	err := collection.Update(bson.M{"email": email}, bson.M{"$set": bson.M{"password": hashedPassword}})

	return err
}

func (u *UserModel) VerifyAccount(email string) error {
	collection := dbConnect.Use(databaseName, "user")

	err := collection.Update(bson.M{"email": email}, bson.M{"$set": bson.M{"is_verified": true}})

	return err
}
