package user

import (
	"belajar-grpc/configs"
	userpb "belajar-grpc/models/user"
	"belajar-grpc/package/connections"
	"database/sql"

	"github.com/sirupsen/logrus"
)

// Server is the server object for this api service.
type Server struct {
	config *configs.Config
	logger *logrus.Logger
	db *sql.DB
	userpb.UnimplementedUsersServer
}

type RegisterUser struct{
	FirstName string `json:"first_name" validate:"required,min=3"`
	LastName string `json:"last_name" validate:"required"`
	Email string `validate:"required"`
	Phone string `validate:"required"`
	Password string `validate:"required"`
}

type User struct {
	Id int64
	FirstName string
	LastName sql.NullString
	Email string
	Phone string
	Role userpb.UserRole
}

func New(config *configs.Config,logger *logrus.Logger) (*Server) {
	db,_:= connections.GetConnection(*config)
	return &Server{config: config, logger: logger,db: db}
}