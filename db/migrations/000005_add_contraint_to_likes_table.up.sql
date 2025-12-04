DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'unique_user_post'
          AND conrelid = 'likes'::regclass
    ) THEN
        ALTER TABLE likes
        ADD CONSTRAINT unique_user_post UNIQUE (user_id, post_id);
    END IF;
END;
$$;