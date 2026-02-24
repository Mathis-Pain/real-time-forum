CREATE TRIGGER IF NOT EXISTS update_post_updated_at_after_comment
AFTER INSERT ON comments
FOR EACH ROW
BEGIN
  UPDATE post
  SET updatedat = CURRENT_TIMESTAMP
  WHERE id = NEW.postid;
END;
