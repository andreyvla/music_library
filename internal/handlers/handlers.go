package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"music_library/internal/models"
	"music_library/internal/service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// key int определяет тип ключа для контекста
type key int

// Ключи для контекста пагинации. Использование типа key предотвращает коллизии
// с другими ключами контекста.
const (
	limitKey key = iota
	offsetKey
)

// Handler содержит сервис для работы с музыкой
type Handler struct {
	musicService service.MusicService
}

// NewHandler создает новый обработчик
func NewHandler(musicService service.MusicService) *Handler {
	return &Handler{musicService: musicService}
}

// CreateSong обрабатывает POST-запрос на создание новой песни
// @Summary Создать песню
// @Description Создает новую песню в библиотеке.
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные песни"
// @Success 201 {object} map[string]int "ID созданной песни"
// @Failure 400 {string} string "Неверный формат JSON"
// @Failure 500 {string} string "Ошибка создания песни"
// @Router /songs [post]
func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		log.Printf("Ошибка декодирования JSON при создании песни: %v", err)
		http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		return
	}

	id, err := h.musicService.AddSong(r.Context(), song)
	if err != nil {
		log.Printf("Ошибка создания песни: %v", err)
		http.Error(w, fmt.Sprintf("ошибка создания песни: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Песня успешно создана, ID: %d", id)

	render.JSON(w, r, map[string]int{"id": id})
}

// GetSongs обрабатывает GET-запрос на получение списка песен.
// @Summary Получить список песен
// @Description Возвращает список песен с пагинацией и фильтрацией.
// @Tags songs
// @Param limit query int false "Количество песен на странице"
// @Param offset query int false "Смещение от начала списка"
// @Param group query string false "Название группы"
// @Param song query string false "Название песни"
// @Param release_date query string false "Дата выпуска"
// @Param link query string false "Ссылка"
// @Success 200 {array} models.Song "Список песен"
// @Failure 500 {string} string "Ошибка получения песен"
// @Router /songs [get]
func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, ok := ctx.Value(limitKey).(int)
	if !ok {
		limit = 10 // Значение по умолчанию, если limit не найден
	}

	offset, ok := ctx.Value(offsetKey).(int)
	if !ok {
		offset = 0 // Значение по умолчанию, если offset не найден
	}

	// Получение фильтров из параметров запроса
	filter := models.Song{
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		ReleaseDate: r.URL.Query().Get("release_date"),
		Link:        r.URL.Query().Get("link"),
	}

	log.Printf("Получение списка песен с limit=%d, offset=%d, filter=%+v", limit, offset, filter)
	songs, err := h.musicService.GetSongs(ctx, limit, offset, filter)
	if err != nil {
		log.Printf("Ошибка получения песен: %v", err)
		http.Error(w, "Ошибка получения песен", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, songs)
}

// GetSong обрабатывает GET-запрос на получение песни по ID.
// @Summary Получить песню по ID
// @Description Возвращает песню по ее ID.
// @Tags songs
// @Param id path int true "ID песни"
// @Success 200 {object} models.Song "Данные песни"
// @Failure 400 {string} string "Неверный ID"
// @Failure 500 {string} string "Ошибка получения песни"
// @Router /songs/{id} [get]
func (h *Handler) GetSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	song, err := h.musicService.GetSongByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("ошибка получения песни: %v", err), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, song)
}

// UpdateSong обрабатывает PUT-запрос на обновление песни.
// @Summary Обновить песню
// @Description Обновляет данные песни.
// @Tags songs
// @Param id path int true "ID песни"
// @Param song body models.Song true "Новые данные песни"
// @Success 200 {object} map[string]string "Статус обновления"
// @Failure 400 {string} string "Неверный ID или формат JSON"
// @Failure 500 {string} string "Ошибка обновления песни"
// @Router /songs/{id} [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "неверный формат JSON", http.StatusBadRequest)
		return
	}

	song.ID = id

	if err := h.musicService.UpdateSong(r.Context(), song); err != nil {
		http.Error(w, fmt.Sprintf("ошибка обновления песни: %v", err), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok"})
}

// DeleteSong обрабатывает DELETE-запрос на удаление песни.
// @Summary Удалить песню
// @Description Удаляет песню по ее ID.
// @Tags songs
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {string} string "Неверный ID"
// @Failure 500 {string} string "Ошибка удаления песни"
// @Router /songs/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	if err := h.musicService.DeleteSong(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf("ошибка удаления песни: %v", err), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{"status": "ok"})
}

// AddVerses добавляет куплеты к песне.
// @Summary Добавить куплеты
// @Description Добавляет куплеты к песне.
// @Tags verses
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param verses body []models.Verse true "Список куплетов"
// @Success 201 {string} string "Куплеты добавлены"
// @Failure 400 {string} string "Неверный ID песни или формат данных"
// @Failure 500 {string} string "Ошибка при добавлении куплетов"
// @Router /songs/{id}/verses [post]
func (h *Handler) AddVerses(w http.ResponseWriter, r *http.Request) {
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Неверный ID песни", http.StatusBadRequest)
		return
	}

	var verses []models.Verse
	if err := json.NewDecoder(r.Body).Decode(&verses); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	for i := range verses {
		verses[i].SongID = songID
	}

	err = h.musicService.AddVerses(r.Context(), songID, verses)
	if err != nil {
		http.Error(w, "Ошибка при добавлении куплетов", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetVerses получает куплеты песни с пагинацией.
// @Summary Получить куплеты
// @Description Возвращает куплеты песни с пагинацией.
// @Tags verses
// @Param id path int true "ID песни"
// @Param limit query int false "Количество куплетов на странице"
// @Param offset query int false "Смещение от начала списка"
// @Success 200 {array} models.Verse "Список куплетов"
// @Failure 400 {string} string "Неверный ID песни или параметры пагинации"
// @Failure 500 {string} string "Ошибка при получении куплетов"
// @Router /songs/{id}/verses [get]
func (h *Handler) GetVerses(w http.ResponseWriter, r *http.Request) {
	songIDStr := chi.URLParam(r, "id")
	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		http.Error(w, "Неверный ID песни", http.StatusBadRequest)
		return
	}

	limit, offset, err := parseLimitOffset(r)
	if err != nil {
		http.Error(w, "Неверные параметры пагинации", http.StatusBadRequest)
		return
	}

	verses, err := h.musicService.GetVerses(r.Context(), songID, limit, offset)
	if err != nil {
		http.Error(w, "Ошибка при получении куплетов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(verses)
}

// Paginate middleware для пагинации.
func Paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit, offset, err := parseLimitOffset(r)
		if err != nil {
			http.Error(w, "Неверные параметры пагинации", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), limitKey, limit)
		ctx = context.WithValue(ctx, offsetKey, offset)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseLimitOffset извлекает параметры пагинации limit и offset из запроса
func parseLimitOffset(r *http.Request) (int, int, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // Значение по умолчанию
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			return 0, 0, fmt.Errorf("неверный параметр limit")
		}
		limit = l
	}

	offset := 0 // Значение по умолчанию
	if offsetStr != "" {
		o, err := strconv.Atoi(offsetStr)
		if err != nil || o < 0 {
			return 0, 0, fmt.Errorf("неверный параметр offset")
		}
		offset = o
	}
	return limit, offset, nil
}

// RootHandler перенаправляет на Swagger документацию
func (h *Handler) RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
}
