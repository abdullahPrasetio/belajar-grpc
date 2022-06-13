package structs

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func StructToStruct(structs1 interface{}, structs2 interface{}) (result string, res bool) {

	m, _ := json.Marshal(structs1)
	var results map[string]interface{}
	json.Unmarshal(m, &results)

	json.Unmarshal([]byte(string(m)), &structs2)

	return string(m), true
}


func ValidationStruct(dataSet interface{}) (bool,map[string][]string){
	validate:=validator.New()
	err:=validate.Struct(dataSet)
    if err != nil {
        //Validation syntax is invalid
        if err,ok := err.(*validator.InvalidValidationError);ok{
            panic(err)
        }

        //Validation errors occurred
        errors := make(map[string][]string)
        //Use reflector to reverse engineer struct
        reflected := reflect.ValueOf(dataSet)
        for _,err := range err.(validator.ValidationErrors){

            // Attempt to find field by name and get json tag name
            field,_ := reflected.Type().FieldByName(err.StructField())
            var name string

            //If json tag doesn't exist, use lower case of name
            if name = field.Tag.Get("json"); name == ""{
                name = strings.ToLower(err.StructField())
            }

            switch err.Tag() {
            case "required":
                errors[name] = append(errors[name], "The "+name+" is required")
                break
            case "email":
                errors[name] = append(errors[name], "The "+name+" should be a valid email")
                break
            case "eqfield":
                errors[name] = append(errors[name], "The "+name+" should be equal to the "+err.Param())
                break
            case "min":
                errors[name] = append(errors[name], "The "+name+" must be at least "+err.Param()+" in length")
                break
            default:
                errors[name] = append(errors[name], "The "+name+" is invalid")
                break
            }
        }

        return false,errors
    }
    return true,nil
}