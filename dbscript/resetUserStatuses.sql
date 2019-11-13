DELIMITER $$
CREATE PROCEDURE resetUserStatuses()
BEGIN
    UPDATE users SET status = 0;
END $$
DELIMITER ;