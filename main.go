package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func connect_db(dsn string) {
	// Connect PostgreSQL
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	log.Print("connected to postgres")
	// Auto migrate in case of schema / model changes
	db.AutoMigrate(&Todo{}, &Task{})
}

func check_env() {
	// check env
	env := os.Getenv("ENV")
	// if env variable is not set load environment from .env file
	if "there" != env {
		// load env
		log.Print("load env")
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

// MAIN
func main() {
	// check env
	check_env()
	// get env
	port := os.Getenv("PORT")
	psql_user := os.Getenv("POSTGRES_USER")
	psql_pw := os.Getenv("POSTGRES_PASSWORD")
	psql_db := os.Getenv("POSTGRES_DB")
	psql_host := os.Getenv("POSTGRES_HOST")
	psql_port := os.Getenv("POSTGRES_PORT")

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", psql_host, psql_user, psql_pw, psql_db, psql_port)
	connect_db(dsn)

	// Start http server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), muxRouter()))

	// golang TLS config
	/*
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		// Create server
		srv := &http.Server{
			Addr:         ":8443",
			Handler:      muxRouter(),
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		// Start https server
		log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))*/

}
