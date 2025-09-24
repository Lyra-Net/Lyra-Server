
CREATE TYPE genre_enum AS
ENUM ('pop', 'rock', 'hiphop', 'jazz', 'classical', 'edm', 'ballad', 'other');

CREATE TYPE mood_enum AS
ENUM ('happy', 'sad', 'chill', 'energetic', 'romantic', 'other');

ALTER TABLE songs ADD COLUMN genre genre_enum;
ALTER TABLE songs ADD COLUMN mood mood_enum;
