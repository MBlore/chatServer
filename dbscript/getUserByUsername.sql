DELIMITER $$
CREATE PROCEDURE getUserByUsername(pUsername varchar(50))
BEGIN
	SELECT id, username, password, displayname, status, statusText FROM users WHERE username = pUsername;
END $$
DELIMITER ;