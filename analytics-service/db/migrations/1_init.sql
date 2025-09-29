-- +goose Up
CREATE TABLE play_events (
    user_id       UUID,
    song_id       String,
    device_id     UUID,
    start_second  UInt32,
    end_second    UInt32,
    timestamp     DateTime DEFAULT now()
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (song_id, timestamp);

CREATE TABLE seek_events (
    user_id      UUID,
    song_id      String,
    device_id    UUID,
    from_second  UInt32,
    to_second    UInt32,
    timestamp    DateTime DEFAULT now()
)
ENGINE = MergeTree
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (song_id, timestamp);


-- +goose Down
DROP TABLE IF EXISTS play_events;
DROP TABLE IF EXISTS seek_events;