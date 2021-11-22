package models

import (
	"golang-base/db"
	"os"
)

var server = os.Getenv("DATABASE")

var databaseName = "golangtodoapi"

var dbConnect = db.NewConnection(server, "golangtodoapi")
