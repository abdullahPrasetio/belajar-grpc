package errors

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Method is the method types.
type Method int
type Methode int
type Service int

type ErrorResponse struct {
	Success bool   `json:"success,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

const (
	FallbackError = `{"status": false, "code": "100", "message": "internal error", "details":""}`
	InternalServerError = `Internal Server Error`

	// List of different Methods
	Authorization Method = iota
	Unauthenticated
	InvalidFormat
	InternalServer
	Unavailable
	Other
)

func CustomHttpError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, ierr error){
	var desc string
	statusProto:=status.Convert(ierr)
	pb := statusProto.Proto()
	w.Header().Set("Content-Type", marshaler.ContentType(pb))
	if s,ok:=status.FromError(ierr);ok {
		w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
		desc=s.Message()
	}else{
		w.WriteHeader(runtime.HTTPStatusFromCode(codes.Unknown))
		desc=ierr.Error()
	}
	b:=new(ErrorResponse)
	err:=json.Unmarshal([]byte(desc),b)
	// log.Println(err)
	if err != nil {
		// log.Println(desc)
		_,_=w.Write([]byte(FallbackError))
	}else{
		err=json.NewEncoder(w).Encode(b)
		if err != nil {
			_,_=w.Write([]byte(FallbackError))
		}
	}
	
}



// FormatError is the exposed function for generating errors.
func FormatError(m Method, c codes.Code, success bool, code string,message string,params interface{}) error {
	// Default
	//errRes := &ErrorResponse{
	//	Success:       	success,
	//	RespCode:       params[0],
	//	RespDesc:       params[1],
	//}
	// Default change to code and msg
	errRes := &ErrorResponse{
		Success: success,
		Code:    code,
		Message:     message,
		Details:	params,
	}

	buf, err := json.Marshal(errRes)

	if err != nil {
		return status.Errorf(c, FallbackError)
	}

	return status.Errorf(c, string(buf))
}

// FormatErrorEncoded is the exposed function for generating errors.
// func FormatErrorArray(w http.ResponseWriter, c codes.Code, message string,params interface{}) error {
// 	var buf []byte
// 	var err error

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(runtime.HTTPStatusFromCode(c))

// 	if len(params) != 2 {
// 		_, _ = w.Write([]byte(InternalServerError))
// 		return status.Errorf(c, InternalServerError)
// 	}

// 	errRes := &ErrorResponse{
// 		Success: success,
// 		Code:    code,
// 		Message:     message,
// 		Details:	params,
// 	}

// 	buf, err = json.Marshal(errRes)
// 	if err != nil {
// 		_, _ = w.Write([]byte(InternalServerError))
// 		return status.Errorf(c, InternalServerError)
// 	}

// 	_, _ = w.Write(buf)

// 	return status.Errorf(c, string(buf))
// }

// FormatError is the exposed function for generating errors.
func FormatErrorHeader(ctx context.Context, md metadata.MD, m Method, c codes.Code, success bool,  code string,message string,params interface{}) error {
	// Default
	//errRes := &ErrorResponse{
	//	Success:       	success,
	//	RespCode:       params[0],
	//	RespDesc:       params[1],
	//}

	// Default change to code and msg
	errRes := &ErrorResponse{
		Success: success,
		Code:    code,
		Message:     message,
		Details:	params,
	}

	// Custom error
	//errResErr := ErrorResponeErr{
	//	Error: ErrorResponse{
	//		Success: success,
	//		Code:    params[0],
	//		Msg:     params[1],
	//	},
	//}

	buf, err := json.Marshal(errRes)

	if err != nil {
		return status.Errorf(c, FallbackError)
	}

	return status.Errorf(c, string(buf))
}

// func to create error from string
func New(errs string) error {
	return errors.New(errs)
}

// WrapErr is the function wrapper to wrap the error
func WrapErr(err error, message string, args ...interface{}) error {
	return errors.Wrapf(err, message, args...)
}
