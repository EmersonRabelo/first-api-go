CREATE OR REPLACE FUNCTION decrement_like_count_on_update()
RETURNS TRIGGER AS $$
BEGIN
  IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
    UPDATE post_likes_count
    SET like_count = like_count - 1
    WHERE post_id = NEW.post_id;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;