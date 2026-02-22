package main

import (
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
	"github.com/parkerjohnson/intercede/services"
	"github.com/parkerjohnson/intercede/web"
	"github.com/supabase-community/supabase-go"
)

func main() {
	err := godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	db, err := sqlx.Connect("pgx", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		log.Fatal("Error connecting to Supabase auth: ", err)
	}
  
	submissionService := services.New(db)
	authService := services.NewAuthService(client)

	wh := webhandlers.NewWeb()
	ah := api.NewApi(submissionService, authService)

	mux := http.NewServeMux()

	//authMiddleWare := middleware.RequireAuth(&middleware.SupabaseTokenVerifier{Client: client})

	mux.HandleFunc("/", wh.Home)
	mux.HandleFunc("/prayerrequest", wh.PrayerRequest)
	mux.HandleFunc("/praisereport", wh.PraiseReport)
	mux.HandleFunc("/adminlogin", wh.AdminLogin)

	mux.HandleFunc("/createPrayer", ah.CreatePrayerRequest)
	mux.HandleFunc("/createPraise", ah.CreatePraiseReport)
	mux.HandleFunc("/adminapilogin", ah.Login)

	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	handler := middleware.Logger(middleware.SecurityHeaders(mux))

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
