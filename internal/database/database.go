package database

import (
	"context"
	"music_library/internal/models"
)

// SongDB - интерфейс для работы с базой данных песен
type SongDB interface {
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

	// GetVersesBySongID получает куплеты песни с пагинацией
	GetVersesBySongID(ctx context.Context, songID, limit, offset int) ([]models.Verse, error)
}
