DELIMITER $$
CREATE PROCEDURE logoffUser(pId int)
BEGIN
    UPDATE users SET status = 0 WHERE id = pId AND status <> 0;
END $$
DELIMITER ;