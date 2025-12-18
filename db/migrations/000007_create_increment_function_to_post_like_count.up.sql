CREATE OR REPLACE FUNCTION increment_like_count()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.deleted_at IS NULL THEN
    INSERT INTO post_likes_count (post_id, like_count)
    VALUES (NEW.post_id, 1)
    ON CONFLICT (post_id) DO
      UPDATE SET like_count = post_likes_count.like_count + 1;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;