-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    username      TEXT    NOT NULL UNIQUE,
    password_hash TEXT    NOT NULL,
    admin         INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS networks (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL,
    address     TEXT    NOT NULL,
    listen_port INTEGER NOT NULL,
    public_key  TEXT    NOT NULL UNIQUE,
    private_key TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS peers (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    network_id    INTEGER NOT NULL REFERENCES networks(id),
    public_key    TEXT    NOT NULL UNIQUE,
    private_key   TEXT    NOT NULL,
    preshared_key TEXT    NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS peer_allowed_ips (
    peer_id INTEGER NOT NULL REFERENCES peers(id),
    ip      TEXT    NOT NULL,
    PRIMARY KEY (peer_id, ip)
);

-- +goose Down
DROP TABLE IF EXISTS peer_allowed_ips;
DROP TABLE IF EXISTS peers;
DROP TABLE IF EXISTS networks;
DROP TABLE IF EXISTS users;
