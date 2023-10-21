package main

import (
	"os"

	"github.com/jfilipedias/codepix-go/application/grpc"
	"github.com/jfilipedias/codepix-go/infra/db"
	"github.com/jinzhu/gorm"
)

var database *gorm.DB

func main() {
	database = db.ConnectDB(os.Getenv("env"))
	grpc.StartGrpcServer(database, 50051)
}
