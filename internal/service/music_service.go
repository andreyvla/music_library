package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"music_library/internal/database"
	"music_library/internal/models"
	"net/http"
	"os"
	"strings"
)

// MusicServiceImpl реализует интерфейс MusicService
type MusicServiceImpl struct {
	db database.SongDB
}

// NewMusicService создает новый MusicServiceImpl
func NewMusicService(db database.SongDB) *MusicServiceImpl {
	return &MusicServiceImpl{db: db}
}

// AddSong добавляет новую песню, предварительно получив информацию из внешнего API
func (s *MusicServiceImpl) AddSong(ctx context.Context, song models.Song) (int, error) {
	log.Printf("Добавление песни: %+v", song)
	details, err := s.GetSongDetails(ctx, song.Group, song.Song)
	if err != nil {
		log.Printf("Ошибка получения деталей песни: %v", err)
		return 0, fmt.Errorf("ошибка получения деталей песни: %w", err)
	}

	song.ReleaseDate = details.ReleaseDate
	song.Link = details.Link

	id, err := s.db.AddSong(ctx, song)
	if err != nil {
		log.Printf("Ошибка добавления песни в БД: %v", err)
		return 0, fmt.Errorf("ошибка добавления песни в БД: %w", err)
	}

	versesText := splitIntoVerses(details.Text)
	verses := make([]models.Verse, len(versesText))
	for i, v := range versesText {
		verses[i] = models.Verse{
			SongID:      id,
			VerseNumber: i + 1,
			Text:        v,
		}
	}

	if err = s.db.AddVerses(ctx, id, verses); err != nil {
		log.Printf("Ошибка добавления куплетов: %v", err)
		return 0, fmt.Errorf("ошибка добавления куплетов: %w", err)
	}

	log.Printf("Песня успешно добавлена. ID: %d", id)
	return id, nil
}

// GetSongDetails получает информацию о песне из внешнего API
func (s *MusicServiceImpl) GetSongDetails(ctx context.Context, group, songTitle string) (models.SongDetails, error) {
	apiUrl := fmt.Sprintf("%s/info?group=%s&song=%s", os.Getenv("API_URL"), group, songTitle)
	log.Printf("Запрос к внешнему API: %s", apiUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		log.Printf("Ошибка создания запроса к API: %v", err)
		return models.SongDetails{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Ошибка выполнения запроса к API: %v", err)
		return models.SongDetails{}, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Неуспешный статус код от API: %d", resp.StatusCode)
		return models.SongDetails{}, fmt.Errorf("неуспешный статус код: %d", resp.StatusCode)
	}

	var details models.SongDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		log.Printf("Ошибка декодирования JSON от API: %v", err)
		return models.SongDetails{}, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	log.Printf("Детали песни получены: %+v", details)
	return details, nil
}

// GetSongs получает список песен из базы данных.
func (s *MusicServiceImpl) GetSongs(ctx context.Context, limit, offset int, filter models.Song) ([]models.Song, error) {
	return s.db.GetSongs(ctx, limit, offset, filter)
}

// GetSongByID получает песню по ID из базы данных.
func (s *MusicServiceImpl) GetSongByID(ctx context.Context, id int) (models.Song, error) {
	return s.db.GetSongByID(ctx, id)
}

// UpdateSong обновляет данные песни в базе данных.
func (s *MusicServiceImpl) UpdateSong(ctx context.Context, song models.Song) error {
	return s.db.UpdateSong(ctx, song)
}

// DeleteSong удаляет песню из базы данных.
func (s *MusicServiceImpl) DeleteSong(ctx context.Context, id int) error {
	return s.db.DeleteSong(ctx, id)
}

func (s *MusicServiceImpl) AddVerses(ctx context.Context, songID int, verses []models.Verse) error {
	return s.db.AddVerses(ctx, songID, verses)
}

func (s *MusicServiceImpl) GetVerses(ctx context.Context, songID, limit, offset int) ([]models.Verse, error) {
	return s.db.GetVersesBySongID(ctx, songID, limit, offset)
}

// splitIntoVerses разбивает текст песни на куплеты по двойному переносу строки.
func splitIntoVerses(text string) []string {
	return strings.Split(text, "\n\n")
}
