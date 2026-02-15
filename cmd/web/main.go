package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/parkerjohnson/intercede/handlers/api"
	"github.com/parkerjohnson/intercede/handlers/webhandlers"
	"github.com/parkerjohnson/intercede/internal/middleware"
	"github.com/parkerjohnson/intercede/web"
)

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	db, err := sqlx.Connect("pgx", os.Getenv("DB_URL"))
	if err != nil {
		fmt.Print(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	wh := webhandlers.NewWeb()
	ah := api.NewApi(db)

	// Create router
	mux := http.NewServeMux()

	mux.HandleFunc("/", wh.Home)
	mux.HandleFunc("/prayerrequest", wh.PrayerRequest)
	mux.HandleFunc("/praisereport", wh.PraiseReport)

	mux.HandleFunc("/createPrayer", ah.CreatePrayerRequest)
	mux.HandleFunc("/createPraise", ah.CreatePraiseReport)

	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	handler := middleware.Logger(middleware.SecurityHeaders(mux))

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
