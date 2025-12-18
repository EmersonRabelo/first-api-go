CREATE UNIQUE INDEX IF NOT EXISTS unique_user_post_active
ON likes(user_id, post_id)
WHERE deleted_at IS NULL;