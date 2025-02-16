package main

import (
	"fmt"
	"log"
	"music_library/internal/database"
	"music_library/internal/handlers"
	"music_library/internal/service"
	"net/http"
	"os"

	_ "music_library/docs" // swag документация для API

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres драйвер
	"github.com/pressly/goose/v3"
)

// @title Music Library API
// @version 1.0
// @description API для работы с музыкальной библиотекой.

// @license.name Unlicensed

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	// Подключение к базе данных
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки подключения к базе данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных!")

	// Применение миграций базы данных
	if err := goose.Up(db.DB, "./migrations"); err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	// Создание сервиса и обработчика
	repo := database.NewPostgresRepository(db)
	musicService := service.NewMusicService(repo)
	handler := handlers.NewHandler(musicService)

	// Создание роутера
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Маршруты
	r.Get("/", handler.RootHandler)
	r.Route("/songs", func(r chi.Router) {
		r.With(handlers.Paginate).Get("/", handler.GetSongs) // GET /songs - получение списка песен
		r.Post("/", handler.CreateSong)                      // POST /songs - создание новой песни
		r.Route("/{id}", func(r chi.Router) {                // Подмаршрутизация для /songs/{id}
			r.Get("/", handler.GetSong)                                 // GET /songs/{id} - получение песни по ID
			r.Put("/", handler.UpdateSong)                              // PUT /songs/{id} - обновление песни
			r.Delete("/", handler.DeleteSong)                           // DELETE /songs/{id} - удаление песни
			r.With(handlers.Paginate).Get("/verses", handler.GetVerses) // GET /songs/{id}/verses - получение куплетов с пагинацией
			r.Post("/verses", handler.AddVerses)                        // POST /songs/{id}/verses - добавление куплетов
		})
	})

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
