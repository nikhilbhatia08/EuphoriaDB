package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/nikhilbhatia08/EuphoriaDB/driver"
	"github.com/nikhilbhatia08/EuphoriaDB/server"
)

func main() {
	dataDir := "./data"
	if len(os.Args) > 1 {
		dataDir = os.Args[1]
	}

	db, err := sql.Open("euphoriadb", dataDir)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	srv := server.NewServer(db, 1000)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("shutting down...")
		srv.Shutdown()
		os.Exit(0)
	}()

	if err := srv.Start(":9909"); err != nil {
		log.Fatal(err)
	}
}
