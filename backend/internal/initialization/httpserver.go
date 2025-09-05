package initialization

import (
	"backend/global"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"go.uber.org/zap"
)

func NewServer(handler http.Handler) *http.Server {
	httpCfg := global.Config.HTTP
	return &http.Server{
		Addr:         ":" + strconv.Itoa(httpCfg.Port),
		Handler:      handler,
		ReadTimeout:  httpCfg.ReadTimeout,
		WriteTimeout: httpCfg.WriteTimeout,
		IdleTimeout:  httpCfg.IdleTimeout,
	}
}

func RunServer(server *http.Server) {
	go func() {
		fmt.Println("Listening and serving HTTP on", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Fatal("Error ListenAndServe():", zap.Error(err))
		}
	}()

	//  kill, Ctrl + C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	// context to grateful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), global.Config.HTTP.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	fmt.Println("Server exited properly")
}
