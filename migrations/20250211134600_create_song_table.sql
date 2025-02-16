-- +goose Up
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    "group" TEXT NOT NULL,
    song TEXT NOT NULL,
    release_date TEXT,
    text TEXT,
    link TEXT
);

-- +goose Down
DROP TABLE songs;
