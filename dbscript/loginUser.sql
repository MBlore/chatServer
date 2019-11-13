DELIMITER $$
CREATE PROCEDURE loginUser(pId int)
BEGIN
    UPDATE users SET status = 1 WHERE id = pId;
END $$
DELIMITER ;