CREATE DATABASE IF NOT EXISTS client;

CREATE TABLE IF NOT EXISTS client.client
(
    id UInt64,
    name String,
    settlement String,
    margin_algorithm UInt8,
    gateway boolean,
    vendor boolean,
    is_active boolean,
    is_pro boolean,
    is_interbank boolean,
    create_at timestamp,
    update_at timestamp
)
ENGINE = MergeTree()
ORDER BY (id)