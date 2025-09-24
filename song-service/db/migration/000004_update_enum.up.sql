BEGIN;

CREATE TYPE genre_enum_new AS ENUM ('Pop', 'Rock', 'HipHop', 'Jazz', 'Classical', 'EDM', 'Ballad', 'Other');

ALTER TABLE songs ALTER COLUMN genre TYPE genre_enum_new
USING
  INITCAP(genre::text)::genre_enum_new;

DROP TYPE genre_enum;

ALTER TYPE genre_enum_new RENAME TO genre_enum;

COMMIT;


BEGIN;

CREATE TYPE mood_enum_new AS ENUM ('Happy', 'Sad', 'Chill', 'Energetic', 'Romantic', 'Other');

ALTER TABLE songs ALTER COLUMN mood TYPE mood_enum_new
USING INITCAP(mood::text)::mood_enum_new;

DROP TYPE mood_enum;

ALTER TYPE mood_enum_new RENAME TO mood_enum;

COMMIT;

-- Migration to update enum values to be capitalized