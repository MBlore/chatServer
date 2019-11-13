DELIMITER $$
CREATE PROCEDURE setStatus(pId int, pStatus int)
BEGIN
    UPDATE users SET status = pStatus WHERE id = pId;
END $$
DELIMITER ;