package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leinonen/bbs/config"
	"github.com/leinonen/bbs/database"
	"github.com/leinonen/bbs/repository"
	"github.com/leinonen/bbs/server"
)

func main() {
	var (
		configFile = flag.String("config", "config.json", "Configuration file path")
		initDB     = flag.Bool("init", false, "Initialize database")
	)
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Printf("Warning: Could not load config file: %v. Using defaults.", err)
		cfg = config.Default()
	}

	db, err := database.Initialize(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repos := repository.NewManager(db)
	defer repos.Close()

	if *initDB {
		if err := database.Migrate(db); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		fmt.Println("Database initialized successfully")
		return
	}

	sshServer := server.NewSSHServer(cfg, repos)

	go func() {
		log.Printf("Starting BBS SSH server on %s", cfg.ListenAddr)
		if err := sshServer.Start(); err != nil {
			log.Fatalf("SSH server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	sshServer.Stop()
}