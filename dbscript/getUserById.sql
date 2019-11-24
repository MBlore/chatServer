DELIMITER $$
CREATE PROCEDURE getUserById(pUserID int)
BEGIN
	SELECT id, username, displayname, status, userimage, statusText FROM users WHERE id = pUserID;
END $$
DELIMITER ;