-- name: GetPendingForUpdate :many
SELECT id, phone_number, content, status, created_at, updated_at, sent_at, response_message_id, fail_reason, retry_count
FROM messages
WHERE status = 'pending'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkProcessing :exec
UPDATE messages
SET status = 'processing'
WHERE id = ANY($1::int[])
  AND status = 'pending';

-- name: MarkSent :exec
UPDATE messages
SET status = 'sent',
    sent_at = now(),
    response_message_id = $2
WHERE id = $1
  AND status IN ('processing','pending');

-- name: MarkFailed :exec
UPDATE messages
SET status = 'failed',
    fail_reason = $2,
    retry_count = retry_count + 1
WHERE id = $1
  AND status IN ('processing','pending');

-- name: ListSent :many
SELECT id, phone_number, content, sent_at, response_message_id
FROM messages
WHERE status = 'sent'
ORDER BY sent_at DESC
LIMIT $1 OFFSET $2;

-- name: InsertMessage :one
INSERT INTO messages (phone_number, content)
VALUES ($1, $2)
RETURNING id, phone_number, content, status, created_at, updated_at, sent_at, response_message_id, fail_reason, retry_count;

-- name: CountPending :one
SELECT COUNT(*) AS cnt FROM messages WHERE status = 'pending';
