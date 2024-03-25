package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"scrapAvito/handler_readiness"
	"scrapAvito/internal/auth"
	"scrapAvito/internal/database"
	"scrapAvito/json_app"
	"scrapAvito/parse/parse_kommersant"
	"scrapAvito/send_info/send_telegram"
	"sync"

	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

type BotTelegram struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	BotToken  string    `json:"bot_token"`
	ChatID    string    `json:"chat_id"`
	UserID    uuid.UUID `json:"user_id"`
}

type SiteParse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	UrlSite   string    `json:"url_site"`
	Type      string    `json:"type"`
}

type SiteParseFollow struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uuid.UUID `json:"user_id"`
	SiteParseID uuid.UUID `json:"site_parse_id"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		APIKey:    dbUser.ApiKey,
	}
}

func databaseBotTelegramToBotTelegram(dbBotTelegram database.Botstelegram) BotTelegram {
	return BotTelegram{
		ID:        dbBotTelegram.ID,
		CreatedAt: dbBotTelegram.CreatedAt,
		UpdatedAt: dbBotTelegram.UpdatedAt,
		Name:      dbBotTelegram.Name,
		BotToken:  dbBotTelegram.BotToken,
		ChatID:    dbBotTelegram.ChatID,
		UserID:    dbBotTelegram.UserID,
	}
}

func databaseSiteParseToSiteParse(dbSiteParse database.Siteparse) SiteParse {
	return SiteParse{
		ID:        dbSiteParse.ID,
		CreatedAt: dbSiteParse.CreatedAt,
		UpdatedAt: dbSiteParse.UpdatedAt,
		Name:      dbSiteParse.Name,
		UrlSite:   dbSiteParse.UrlSite,
		Type:      dbSiteParse.Type,
	}
}

func databaseSiteParseFollowToSiteParseFollow(dbSiteParseFollow database.Siteparsefollow) SiteParseFollow {
	return SiteParseFollow{
		ID:          dbSiteParseFollow.ID,
		CreatedAt:   dbSiteParseFollow.CreatedAt,
		UpdatedAt:   dbSiteParseFollow.UpdatedAt,
		UserID:      dbSiteParseFollow.UserID,
		SiteParseID: dbSiteParseFollow.SiteParseID,
	}
}

func main() {

	// parse_kommersant.ParseKommersant()

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found")
	}

	log.Printf("Server starting on port %s", portString)

	//установка БД
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	db := database.New(conn)

	apiCfg := apiConfig{
		DB: db,
	}
	//Конец установки БД

	go startScraping(db, 10, time.Minute)

	// data, err := apiCfg.DB.GetSiteParseFollows(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("SiteParseID: %s\n", data[0].SiteParseID)
	// fmt.Printf("UserID: %s\n", data[0].UserID)

	// data1, err := apiCfg.DB.SelectSiteParse(context.Background(), data[0].SiteParseID)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// dataParse := parse_kommersant.ParseKommersant(data1.UrlSite)
	// text := ""
	// for _, item := range dataParse.Items {

	// 	_, err := db.CreateNewsElement(context.Background(), database.CreateNewsElementParams{
	// 		ID:          uuid.New(),
	// 		CreatedAt:   time.Now().UTC(),
	// 		UpdatedAt:   time.Now().UTC(),
	// 		Title:       item.Title,
	// 		NewsDate:    item.DateTime,
	// 		Url:         item.Url,
	// 		SiteParseID: data[0].SiteParseID,
	// 	})
	// 	if err != nil {
	// 		if err != sql.ErrNoRows {
	// 			log.Println("failed to create news element", err)
	// 		}
	// 	} else {
	// 		text = fmt.Sprintf("%sДата и время: %s\nНазвание: %s\nСсылка: %s\n\n", text, item.DateTime, item.Title, item.Url)
	// 		fmt.Printf("Дата и время: %s\nНазвание: %s\nСсылка: %s\n\n", item.DateTime, item.Title, item.Url)
	// 	}
	// }

	// encodedText := url.QueryEscape(text)
	// data2, err := apiCfg.DB.SelectBotTelegram(context.Background(), data[0].UserID)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if text != "" {
	// 	send_telegram.SendTelegram(data2.BotToken, data2.ChatID, encodedText)
	// }

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handler_readiness.HandlerReadiness)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/botTelegram", apiCfg.middlewareAuth(apiCfg.handlerCreateBotTelegram))
	v1Router.Post("/siteParse", apiCfg.handlerCreateSiteParse)
	v1Router.Post("/siteParseFollow", apiCfg.middlewareAuth(apiCfg.handlerCreateSiteParseFollows))
	v1Router.Get("/clearNewsElement", apiCfg.handlerClearNewsElement)
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Port: %s\n", portString)

}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Couldn't create user:%v", err))
		return
	}
	json_app.RespondWithJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	json_app.RespondWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerCreateBotTelegram(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name     string `json:"name"`
		BotToken string `json:"botToken"`
		ChatID   string `json:"chatId"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	botTelegram, err := apiCfg.DB.CreateBotTelegram(r.Context(), database.CreateBotTelegramParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		BotToken:  params.BotToken,
		ChatID:    params.ChatID,
		UserID:    user.ID,
	})
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Couldn't create botTelegram:%v", err))
		return
	}
	json_app.RespondWithJSON(w, 201, databaseBotTelegramToBotTelegram(botTelegram))
}

func (apiCfg *apiConfig) handlerCreateSiteParse(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name    string `json:"name"`
		UrlSite string `json:"urlSite"`
		Type    string `json:"type"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	siteParse, err := apiCfg.DB.CreateSiteParse(r.Context(), database.CreateSiteParseParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		UrlSite:   params.UrlSite,
		Type:      params.Type,
	})
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Couldn't create siteParse:%v", err))
		return
	}
	json_app.RespondWithJSON(w, 201, databaseSiteParseToSiteParse(siteParse))
}

func (apiCfg *apiConfig) handlerCreateSiteParseFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		SiteParseId uuid.UUID `json:"siteParseId"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}

	siteParseFollow, err := apiCfg.DB.CreateSiteParseFollows(r.Context(), database.CreateSiteParseFollowsParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		SiteParseID: params.SiteParseId,
		UserID:      user.ID,
	})
	if err != nil {
		json_app.RespondWithError(w, 400, fmt.Sprintf("Couldn't create siteParseFollow:%v", err))
		return
	}
	json_app.RespondWithJSON(w, 201, databaseSiteParseFollowToSiteParseFollow(siteParseFollow))
}

func (apiCfg *apiConfig) handlerClearNewsElement(w http.ResponseWriter, r *http.Request) {
	apiCfg.DB.ClearNewsElement(r.Context())
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			json_app.RespondWithError(w, 403, fmt.Sprintf("Auth error:%v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			json_app.RespondWithError(w, 400, fmt.Sprintf("Couldn't get user:%v", err))
			return
		}

		handler(w, r, user)
	}
}

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		log.Println("Tick")
		siteParses, err := db.GetNextSiteParseToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching site parse:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, siteParse := range siteParses {
			wg.Add(1)

			go scrapeSiteParse(db, wg, siteParse)
		}
		wg.Wait()
	}
}

func scrapeSiteParse(db *database.Queries, wg *sync.WaitGroup, siteParse database.Siteparse) {
	defer wg.Done()
	log.Println("ParseKommersant")
	dataParse := parse_kommersant.ParseKommersant(siteParse.UrlSite)
	text := ""
	for _, item := range dataParse.Items {

		_, err := db.CreateNewsElement(context.Background(), database.CreateNewsElementParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			NewsDate:    item.DateTime,
			Url:         item.Url,
			SiteParseID: siteParse.ID,
		})
		if err != nil {
			if err != sql.ErrNoRows {
				log.Println("failed to create news element", err)
			}
		} else {
			text = fmt.Sprintf("%sДата и время: %s\nНазвание: %s\nСсылка: %s\n\n", text, item.DateTime, item.Title, item.Url)
			fmt.Printf("Дата и время: %s\nНазвание: %s\nСсылка: %s\n\n", item.DateTime, item.Title, item.Url)
		}
	}

	encodedText := url.QueryEscape(text)

	internalWg := &sync.WaitGroup{}
	botData, err := db.GetBotDataBySiteParseID(context.Background(), siteParse.ID)
	if err != nil {
		log.Println("error get bot data:", err)
	}

	internalWg.Add(len(botData))
	for _, telegramUser := range botData {
		go sendDataToTelegram(telegramUser, internalWg, encodedText)
	}
	internalWg.Wait()

}

func sendDataToTelegram(telegramUser database.GetBotDataBySiteParseIDRow, internalWg *sync.WaitGroup, encodedText string) {
	defer internalWg.Done()
	send_telegram.SendTelegram(telegramUser.BotToken, telegramUser.ChatID, encodedText)
	// log.Println(telegramUser.BotToken)
	// log.Println(telegramUser.ChatID)
}
