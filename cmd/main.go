package main

import (
	"context"
	"encoding/json"
	"fmt"
	workerpool "gds/worker-pool"
	"io"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type ReqBody struct {
	Num int
}

var WorkerPool *workerpool.WorkerPool

func main() {

	server := &http.Server{
		Addr: ":8080",
	}

	protected := authMiddleware(adminHandler)
	http.HandleFunc("/start-job", protected)
	WorkerPool = workerpool.InitWorker(2)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	fmt.Println("Shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown forcefully: %v", err)
	}
	WorkerPool.Shutdown()

	log.Println("Server exited cleanly")

}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleBasicAuth(next, w, r)
	}
}

func handleBasicAuth(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()

	if !ok || user != "admin" || pass != "password123" {
		w.Header().Set("WWW-Authenticate", `Basic realm = Restricted`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	next(w, r)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var req ReqBody
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		fmt.Fprintf(w, "error decoding object")
		return
	}

	WorkerPool.Submit(req.Num)
	fmt.Fprintf(w, "job accepted: %d", req.Num)
}
