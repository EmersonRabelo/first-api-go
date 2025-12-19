CREATE TRIGGER trg_increment_reply
AFTER INSERT ON replies
FOR EACH ROW
EXECUTE FUNCTION increment_reply_count();