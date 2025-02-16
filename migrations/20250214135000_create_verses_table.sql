-- +goose Up
CREATE TABLE verses (
    id SERIAL PRIMARY KEY,
    song_id INTEGER REFERENCES songs(id) ON DELETE CASCADE,
    verse_number INTEGER NOT NULL,
    text TEXT NOT NULL
);

-- +goose Down
DROP TABLE verses;
