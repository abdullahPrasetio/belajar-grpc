package router

import (
	"belajar-grpc/configs"
	"belajar-grpc/package/errors"
	"belajar-grpc/package/header"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	userpb "belajar-grpc/models/user"
)

func NewHTTPServer(config *configs.Config, logger *logrus.Logger) error {
	

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect to the GRPC server
	addr := "0.0.0.0:" + config.PORT_GRPC
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return err
	}
	defer conn.Close()
	// Create new grpc-gateway
	rmux:=gwruntime.NewServeMux(gwruntime.WithErrorHandler(errors.CustomHttpError),gwruntime.WithForwardResponseOption(header.HttpResponseModifier),
		gwruntime.WithIncomingHeaderMatcher(CustomMatcher),
	)

	// register gateway endpoints
	for _,f :=range []func(ctx context.Context,mux *gwruntime.ServeMux,conn *grpc.ClientConn)error {
		userpb.RegisterUsersHandler,
	}{
		if err=f(ctx,rmux,conn);err!=nil {
			log.Fatal(err)
			return err
		}
	}
	// rmux.HandlePath("GET","/api/v1/users", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	// 	w.Write([]byte("hello "))
	// })
	mux:=http.NewServeMux()
	mux.Handle("/",rmux)
		
	// run swagger server
	if config.APP_ENV == "local" {
		CreateSwagger(mux)
	}
	// server:=http.Server{Handler:withLogger(mux),}

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "X-CSRF-Token", "Authorization", "Timezone-Offset"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// running rest http server
	err = http.ListenAndServe("0.0.0.0:"+config.PORT, handlers.CORS(headersOk, originsOk, methodsOk)(mux))
	log.Println("[SERVER] REST HTTP server is ready")
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
	
}


func CustomMatcher(key string) (string, bool) {
	switch key {
		case "X-User-Id":
		return key, true
		case "X-Api-Key":
		return key, true
		default:
		return gwruntime.DefaultHeaderMatcher(key)
	}
   }


func CreateSwagger(gwmux *http.ServeMux) {
	// register swagger service server
	gwmux.HandleFunc("/corp/rest/v1/user/docs/user.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "swagger/proto/user.swagger.json")
	})

	// load swagger-ui file
	fs := http.FileServer(http.Dir("swagger/swagger-ui"))
	gwmux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))
}