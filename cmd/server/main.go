package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/IgorChicherin/gophkeeper/internal/app/server/http/router"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/authlib/sha256"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/config"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/crypto/crypto509"
	"github.com/IgorChicherin/gophkeeper/internal/pkg/db"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.GetServerConfig()

	certsManager := crypto509.NewCertsManager(cfg.PrivateCertPath, cfg.PublicCertPath)

	private, public, err := certsManager.GetCerts()

	if err != nil {
		log.Fatalf("unable to load certs: %s", err)
	}

	if err != nil {
		log.Fatalf("unable to config server: %s", err)
	}
	ctxDB := context.Background()
	conn, err := pgx.Connect(ctxDB, cfg.Database)

	if err != nil {
		log.Fatalf("unable to connect DB: %s", err)
	}

	if err := db.Migrate(cfg.Database); err != nil {
		log.Fatalf("migration failed: %s", err)
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.WithFields(log.Fields{"func": "main"}).Errorln(err)
		}
	}(conn, ctxDB)

	err = conn.Ping(ctxDB)

	if err != nil {
		log.Fatalf("unable to connect DB: %s", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	hashService := sha256.NewSha256HashService(cfg.Key)
	decrypter, err := crypto509.NewDecrypter(private)

	if err != nil {
		log.Fatalf("unable to create encrypting service: %s", err)
	}

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router.NewRouter(conn, hashService, public, decrypter),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Infoln("Server Started")
	<-done
	log.Infoln("Server Stopped")

	ctx := context.Background()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Infoln("Server Exited Properly")
}
