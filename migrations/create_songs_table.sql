CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    music_group VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    releaseDate VARCHAR(50) NOT NULL,
    text TEXT NOT NULL,
    link VARCHAR(255) NOT NULL,
    UNIQUE (music_group, song)
);

