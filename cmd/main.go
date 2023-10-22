package main

import (
	"github.com/jfilipedias/codepix-go/application/grpc"
	"github.com/jfilipedias/codepix-go/infra/db"
	"gorm.io/gorm"
)

var database *gorm.DB

func main() {
	database = db.ConnectDB()
	grpc.StartGrpcServer(database, 50051)
}
