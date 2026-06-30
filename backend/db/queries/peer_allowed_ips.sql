-- name: InsertPeerAllowedIP :exec
INSERT INTO peer_allowed_ips (peer_id, ip) VALUES (?, ?);

-- name: GetPeerAllowedIPs :many
SELECT ip FROM peer_allowed_ips WHERE peer_id = ?;

-- name: GetAllowedIPsByNetworkID :many
SELECT peer_id, ip FROM peer_allowed_ips
WHERE peer_id IN (SELECT id FROM peers WHERE network_id = ?);

-- name: DeletePeerAllowedIPs :exec
DELETE FROM peer_allowed_ips WHERE peer_id = ?;

-- name: DeleteAllowedIPsByNetworkID :exec
DELETE FROM peer_allowed_ips WHERE peer_id IN (
    SELECT id FROM peers WHERE network_id = ?
);
