CREATE TRIGGER trg_decrement_reply
AFTER UPDATE ON replies
FOR EACH ROW
EXECUTE FUNCTION decrement_reply_count_on_update();