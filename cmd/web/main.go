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
	stripe "github.com/stripe/stripe-go/v82"
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

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	stripePriceID := os.Getenv("STRIPE_PRICE_ID")
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:" + port
	}

	submissionService := services.New(db)
	authService := services.NewAuthService(client)
	churchService := services.NewChurchService(db)

	wh := webhandlers.NewWeb(churchService, submissionService)
	ah := api.NewApi(submissionService, authService)
	ch := api.NewChurchAPI(churchService, nil, stripeWebhookSecret, stripePriceID, baseURL, client.Auth)

	mux := http.NewServeMux()

	authMiddleWare := middleware.RequireAuth(&middleware.SupabaseTokenVerifier{Client: client})

	mux.HandleFunc("/", wh.Home)
	mux.HandleFunc("/prayerrequest", wh.PrayerRequest)
	mux.HandleFunc("/praisereport", wh.PraiseReport)
	mux.HandleFunc("/adminlogin", wh.AdminLogin)
	mux.HandleFunc("/church-create", wh.ChurchCreate)
	mux.HandleFunc("/church/success", wh.ChurchSuccess)
	mux.HandleFunc("/church/code-status", wh.ChurchCodeStatus)
	mux.HandleFunc("/church/cancel", wh.ChurchCancel)
	mux.Handle("/admin/dashboard", authMiddleWare(http.HandlerFunc(wh.AdminDashboard)))

	mux.HandleFunc("/createPrayer", ah.CreatePrayerRequest)
	mux.HandleFunc("/createPraise", ah.CreatePraiseReport)
	mux.HandleFunc("/adminapilogin", ah.Login)
	mux.HandleFunc("/api/church/checkout", ch.StartCheckout)
	mux.HandleFunc("/webhooks/stripe", ch.HandleStripeWebhook)
	mux.Handle("/admin/dashboard/submissions", authMiddleWare(http.HandlerFunc(wh.DashboardSubmissions)))

	staticFS, err := fs.Sub(web.StaticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	handler := middleware.Logger(middleware.SecurityHeaders(mux))

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
