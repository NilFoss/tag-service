package main

import (
	"context"
	"flag"
	pb "github.com/go-programming-tour-book/tag-service/proto"
	"github.com/go-programming-tour-book/tag-service/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net/http"
	"strings"
)

//func main()  {
//	s := grpc.NewServer()
//	pb.RegisterTagServiceServer(s,server.NewTagServer())
//
//	reflection.Register(s)
//
//	lis,err := net.Listen(`tcp`,`:8080`)
//	if err != nil {
//		log.Fatalf("dasdada: %v",err)
//	}
//
//	err = s.Serve(lis)
//
//	if err != nil {
//		log.Fatalf(`asdsdsad: %v`,err)
//	}
//
//
//}

var port string

func init() {
	flag.StringVar(&port, "port", "8004", "test")
}

func main() {
	err := RunServer("8004")
	if err != nil {
		log.Fatalf("ssss %v", err)
	}
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func RunServer(port string) error {
	httpMux := runHttpServer()
	gatewayMux := runGrpcGatewayServer()
	httpMux.Handle("/", gatewayMux)

	return http.ListenAndServe(":"+port, grpcHandlerFunc(runGrpcServer(), httpMux))
}

func runHttpServer() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`pong`))
	})

	return serveMux
}

func runGrpcServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTagServiceServer(s, server.NewTagServer())
	reflection.Register(s)

	return s
}

func runGrpcGatewayServer() *runtime.ServeMux {
	endpoint := "0.0.0.0:" + port
	gwmux := runtime.NewServeMux()
	dopts := []grpc.DialOption{grpc.WithInsecure()}

	_ = pb.RegisterTagServiceHandlerFromEndpoint(context.Background(), gwmux, endpoint, dopts)
	return gwmux
}
