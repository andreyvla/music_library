package service

import (
	"context"
	"music_library/internal/models"
)

// MusicService описывает интерфейс сервиса для работы с музыкальной библиотекой
type MusicService interface {
	// AddSong добавляет новую песню
	AddSong(ctx context.Context, song models.Song) (int, error)

	// GetSongs получает список песен с фильтрацией и пагинацией
	GetSongs(ctx context.Context, limit, offset int, filter models.Song) ([]models.Song, error)

	// GetSongByID получает песню по ID
	GetSongByID(ctx context.Context, id int) (models.Song, error)

	// UpdateSong обновляет данные песни
	UpdateSong(ctx context.Context, song models.Song) error

	// DeleteSong удаляет песню по ID
	DeleteSong(ctx context.Context, id int) error

	// AddVerses добавляет куплеты к песне
	AddVerses(ctx context.Context, songID int, verses []models.Verse) error

	// GetVerses получает куплеты песни с пагинацией
	GetVerses(ctx context.Context, songID, limit, offset int) ([]models.Verse, error)

	// GetSongDetails получает информацию о песне из внешнего API
	GetSongDetails(ctx context.Context, group, song string) (models.SongDetails, error)
}
