package user

import (
	userpb "belajar-grpc/models/user"
	"belajar-grpc/package/errors"
	"belajar-grpc/package/structs"
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

func (s *Server) Register(ctx context.Context, param *userpb.UserRegister)(*userpb.ResponseUserData,error){
	userRegister:=RegisterUser{}
	// respon:=ResponseMessage{
	// 	Data:RegisterUser{},
	// }
	// userRegister=RegisterUser(&param)
	structs.StructToStruct(param, &userRegister)
	// validate:=validator.New()
	// rules := map[string]string{
	// 	"FirstName": "min=4,max=6",
	// 	"LastName":  "min=4,max=6",
	// }
	// validate.RegisterStructValidationMapRules(rules, param)
	// err := validate.Struct(param)
	// fmt.Println(err)
	// fmt.Println()
	valid,respo:=structs.ValidationStruct(userRegister)
	if !valid {
		return nil,errors.FormatError(errors.InternalServer,codes.InvalidArgument,false,"400","Validation Error",respo)
	}
	// err:=param.ValidateAll()
	
	// if err!=nil {
	// 	log.Println(err)
	// 	data :=map[string][]string{
	// 		"name":{"error name is required"},
	// 		"phone":{"phone is required"},
	// 	}
	// 	return nil,errors.FormatError(errors.InternalServer,codes.Internal,false,"400","Validation Error",data)
	// }

	password :=[]byte(param.Password)
	hashedPassword,err :=bcrypt.GenerateFromPassword(password,bcrypt.DefaultCost)
	if err != nil{
		return nil,err
	}
	sql :="INSERT INTO users (first_name,last_name,email,password,phone,role) VALUES (?,?,?,?,?,?)"
	res,err:=s.db.ExecContext(ctx,sql,param.FirstName,param.LastName,param.Email,hashedPassword,param.Phone,"0")
	
	if err!=nil {
		log.Println("Error inserting user",err.Error())
	}
	id, err := res.LastInsertId()
    if err != nil {
        return nil,err
    }
	data:=[]*userpb.UserWithoutPassword{
		{
			Id:id,
			FirstName:param.FirstName,
			LastName:param.LastName,
			Email:param.Email,
			Phone:param.Phone,
		},
	}
	
	return &userpb.ResponseUserData{
		Message:"Registered Success",
		Status: "Success",
		Data: data,
	},nil
}

