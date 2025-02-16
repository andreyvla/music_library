package database

import (
	"context"
	"fmt"
	"log"
	"music_library/internal/models"

	"github.com/jmoiron/sqlx"
)

// PostgresRepository реализует интерфейс SongDB для PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository создает новый PostgresRepository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// AddSong добавляет новую песню в базу данных
func (r *PostgresRepository) AddSong(ctx context.Context, song models.Song) (int, error) {
	query := `
		INSERT INTO songs ("group", song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowxContext(ctx, query, song.Group, song.Song, song.ReleaseDate, song.Link).Scan(&id)
	if err != nil {
		log.Printf("Ошибка добавления песни: %v", err)
		return 0, fmt.Errorf("ошибка добавления песни: %w", err)
	}

	log.Printf("Песня добавлена, ID: %d", id)
	return id, nil
}

// GetSongs получает список песен из базы данных с учетом фильтрации и пагинации
func (r *PostgresRepository) GetSongs(ctx context.Context, limit, offset int, filter models.Song) ([]models.Song, error) {
	// Базовый запрос
	query := `
        SELECT id, "group", song, release_date, text, link
        FROM songs
        WHERE 1=1
    `

	// Параметры для запроса
	args := []interface{}{}
	argIndex := 1

	// Добавление условий фильтрации
	if filter.Group != "" {
		query += fmt.Sprintf(` AND "group" ILIKE $%d`, argIndex)
		args = append(args, "%"+filter.Group+"%")
		argIndex++
	}
	if filter.Song != "" {
		query += fmt.Sprintf(` AND song ILIKE $%d`, argIndex)
		args = append(args, "%"+filter.Song+"%")
		argIndex++
	}
	if filter.ReleaseDate != "" {
		query += fmt.Sprintf(` AND release_date = $%d`, argIndex)
		args = append(args, filter.ReleaseDate)
		argIndex++
	}

	if filter.Link != "" {
		query += fmt.Sprintf(` AND link = $%d`, argIndex)
		args = append(args, filter.Link)
		argIndex++
	}

	// Добавление пагинации
	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	var songs []models.Song
	err := r.db.SelectContext(ctx, &songs, query, args...)
	if err != nil {
		log.Printf("Ошибка получения песен: %v", err)
		return nil, fmt.Errorf("ошибка получения песен: %w", err)
	}

	return songs, nil
}

// GetSongByID получает песню по ID
func (r *PostgresRepository) GetSongByID(ctx context.Context, id int) (models.Song, error) {
	query := `
		SELECT id, "group", song, release_date, text, link
		FROM songs
		WHERE id = $1
	`
	var song models.Song
	err := r.db.GetContext(ctx, &song, query, id)
	if err != nil {
		log.Printf("Ошибка получения песни по ID: %v", err)
		return models.Song{}, fmt.Errorf("ошибка получения песни по ID: %w", err)
	}

	return song, nil
}

// UpdateSong обновляет данные песни
func (r *PostgresRepository) UpdateSong(ctx context.Context, song models.Song) error {
	query := `
		UPDATE songs
		SET "group" = $1, song = $2, release_date = $3, text = $4, link = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query, song.Group, song.Song, song.ReleaseDate, song.Link, song.ID)
	if err != nil {
		log.Printf("Ошибка обновления песни: %v", err)
		return fmt.Errorf("ошибка обновления песни: %w", err)
	}

	log.Printf("Песня обновлена, ID: %d", song.ID)
	return nil
}

// DeleteSong удаляет песню из базы данных
func (r *PostgresRepository) DeleteSong(ctx context.Context, id int) error {
	query := `
		DELETE FROM songs
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("Ошибка удаления песни: %v", err)
		return fmt.Errorf("ошибка удаления песни: %w", err)
	}

	log.Printf("Песня удалена, ID: %d", id)

	return nil
}

// AddVerses добавляет куплеты для песни
func (r *PostgresRepository) AddVerses(ctx context.Context, songID int, verses []models.Verse) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Ошибка начала транзакции: %v", err)
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	for _, verse := range verses {
		query := `
			INSERT INTO verses (song_id, verse_number, text)
			VALUES ($1, $2, $3)
		`
		_, err := tx.ExecContext(ctx, query, songID, verse.VerseNumber, verse.Text)
		if err != nil {
			log.Printf("Ошибка добавления куплета: %v", err)
			return fmt.Errorf("ошибка добавления куплета: %w", err)
		}
		log.Printf("Куплет добавлен, SongID: %d, VerseNumber: %d", songID, verse.VerseNumber)
	}

	return tx.Commit()
}

// GetVersesBySongID получает куплеты для песни с пагинацией
func (r *PostgresRepository) GetVersesBySongID(ctx context.Context, songID, limit, offset int) ([]models.Verse, error) {
	query := `
		SELECT id, song_id, verse_number, text
		FROM verses
		WHERE song_id = $1
		LIMIT $2 OFFSET $3
	`

	var verses []models.Verse
	err := r.db.SelectContext(ctx, &verses, query, songID, limit, offset)
	if err != nil {
		log.Printf("Ошибка получения куплетов: %v", err)
		return nil, fmt.Errorf("ошибка получения куплетов: %w", err)
	}

	return verses, nil
}
