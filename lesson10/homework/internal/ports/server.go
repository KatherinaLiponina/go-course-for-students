package ports

import (
	"context"
	"errors"
	"fmt"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app"
	grpc_func "homework10/internal/ports/grpc"
	"homework10/internal/ports/httpgin"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"golang.org/x/sync/errgroup"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	api := handler.Group("/api/v1")
	httpgin.AppRouter(api, a)
	s := &http.Server{Addr: port, Handler: handler}
	return s
}

func NewGRPCServer(port string, a app.App) *grpc.Server {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}
	service := &grpc_func.AdUserService{App: a}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(UnaryServerInterceptor),
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor(opts...)))
	grpc_func.RegisterAdServiceServer(server, service)
	return server
}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println(time.Now().GoString() + ": " + info.FullMethod)

	return handler(ctx, req)
}

const (
	grpcPort = ":50054"
	httpPort = ":18080"
)

func CreateServer(ctx context.Context, ch chan int) (*http.Server, *grpc.Server) {

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	a := app.NewApp(adrepo.New(), userrepo.New())
	httpServer := NewHTTPServer(httpPort, a)
	grpcServer := NewGRPCServer(grpcPort, a)

	eg, ctx := errgroup.WithContext(ctx)

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		eg.Go(func() error {
			select {
			case s := <-sigQuit:
				log.Printf("captured signal: %v\n", s)
				return fmt.Errorf("captured signal: %v", s)
			case <-ctx.Done():
				return nil
			}
		})

		// run grpc server
		eg.Go(func() error {
			log.Printf("starting grpc server, listening on %s\n", grpcPort)
			defer log.Printf("close grpc server listening on %s\n", grpcPort)

			errCh := make(chan error)

			defer func() {
				grpcServer.GracefulStop()
				_ = lis.Close()

				close(errCh)
			}()

			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					errCh <- err
				}
			}()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errCh:
				return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
			}
		})

		eg.Go(func() error {
			log.Printf("starting http server, listening on %s\n", httpServer.Addr)
			defer log.Printf("close http server listening on %s\n", httpServer.Addr)

			errCh := make(chan error)

			defer func() {
				shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := httpServer.Shutdown(shCtx); err != nil {
					log.Printf("can't close http server listening on %s: %s", httpServer.Addr, err.Error())
				}

				close(errCh)
			}()

			go func() {
				if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
					errCh <- err
				}
			}()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errCh:
				return fmt.Errorf("http server can't listen and serve requests: %w", err)
			}
		})

		if err := eg.Wait(); err != nil {
			log.Printf("gracefully shutting down the servers: %s\n", err.Error())
		}

		log.Println("servers were successfully shutdown")

		ch <- 0
	}()
	time.Sleep(time.Millisecond * 30)
	return httpServer, grpcServer
}
