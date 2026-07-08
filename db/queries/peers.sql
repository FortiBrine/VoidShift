-- name: CreatePeer :one
INSERT INTO peers (network_id, public_key, private_key, preshared_key)
VALUES (?, ?, ?, ?)
RETURNING id, network_id, public_key, private_key, preshared_key;

-- name: GetPeer :one
SELECT id, network_id, public_key, private_key, preshared_key
FROM peers WHERE id = ?;

-- name: GetPeersByNetworkID :many
SELECT id, network_id, public_key, private_key, preshared_key
FROM peers WHERE network_id = ?;

-- name: DeletePeer :exec
DELETE FROM peers WHERE id = ?;

-- name: DeletePeersByNetworkID :exec
DELETE FROM peers WHERE network_id = ?;

-- name: CountPeers :one
SELECT COUNT(*) FROM peers;
