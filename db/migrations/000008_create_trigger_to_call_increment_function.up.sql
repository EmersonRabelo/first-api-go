CREATE TRIGGER trg_increment_like
AFTER INSERT ON likes
FOR EACH ROW
EXECUTE FUNCTION increment_like_count();