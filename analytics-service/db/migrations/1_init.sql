-- +goose Up
CREATE TABLE play_events (
    user_id       String,
    song_id       String,
    device_id     String,
    start_second  UInt32,
    end_second    UInt32,
    timestamp     DateTime DEFAULT now()
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (song_id, timestamp);


-- +goose Down
DROP TABLE play_events;