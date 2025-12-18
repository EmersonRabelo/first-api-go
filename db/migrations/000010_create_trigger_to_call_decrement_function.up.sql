CREATE TRIGGER trg_decrement_like
AFTER UPDATE ON likes
FOR EACH ROW
EXECUTE FUNCTION decrement_like_count_on_update();