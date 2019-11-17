DELIMITER $$
CREATE PROCEDURE getUserContactByContactUserID(pUserID int, pContactUserID int)
BEGIN
	SELECT id FROM friends WHERE userid = pUserID AND frienduserid = pContactUserID;
END $$
DELIMITER ;