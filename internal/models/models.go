package models

// Song представляет песню в музыкальной библиотеке
type Song struct {
	ID          int      `db:"id" json:"id"`
	Group       string   `db:"group" json:"group"`
	Song        string   `db:"song" json:"song"`
	ReleaseDate string   `db:"release_date" json:"release_date"`
	Link        string   `db:"link" json:"link"`
	Verses      []*Verse `db:"-" json:"verses"`
}

// SongDetails представляет информацию о песне из внешнего API
type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Verse представляет куплет песни
type Verse struct {
	ID          int    `db:"id" json:"id"`
	SongID      int    `db:"song_id" json:"song_id"`
	VerseNumber int    `db:"verse_number" json:"verse_number"`
	Text        string `db:"text" json:"text"`
}
