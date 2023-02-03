DELIMITER $$
CREATE FUNCTION get_video_comment_cnt(id int)
    RETURNS int
READS SQL DATA
DETERMINISTIC
BEGIN
    DECLARE cnt int;
    SELECT COUNT(video_id) FROM favourite_relations WHERE video_id=id INTO cnt;
    RETURN cnt;
END