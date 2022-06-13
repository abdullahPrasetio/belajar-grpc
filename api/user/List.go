package user

import (
	userpb "belajar-grpc/models/user"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)


func (s *Server) List(ctx context.Context,void *empty.Empty)(*userpb.ResponseUserData, error){
	sql:="SELECT id,first_name,last_name,email,phone,role FROM users"
	response:=userpb.ResponseUserData{
		Message: "Success get List User",
		Status: "Success",
	}

	rows,err:=s.db.QueryContext(ctx,sql)
	if err != nil {
		return &response,err
	}
	defer rows.Close()
	for rows.Next() {
		user:=User{}
		rows.Scan(&user.Id,&user.FirstName,&user.LastName,&user.Email,&user.Phone,&user.Role)
		userData:=formatToUserWithoutPassword(user)
		response.Data = append(response.Data, &userData)
	}
	
	return &response, nil
}

func formatToUserWithoutPassword(user User) (users userpb.UserWithoutPassword){
	users.Id=user.Id
	users.FirstName=user.FirstName
	users.LastName=user.LastName.String
	users.Email=user.Email
	users.Phone=user.Phone
	users.Role=user.Role

	return users
}
