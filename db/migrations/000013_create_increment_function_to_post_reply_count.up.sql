CREATE OR REPLACE FUNCTION increment_reply_count()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.deleted_at IS NULL THEN
    INSERT INTO post_replies_count (post_id, reply_count)
    VALUES (NEW.post_id, 1)
    ON CONFLICT (post_id) DO
      UPDATE SET reply_count = post_replies_count.reply_count + 1;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;