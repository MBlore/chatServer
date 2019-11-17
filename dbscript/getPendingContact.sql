DELIMITER $$
CREATE PROCEDURE getPendingContact(pUserID int, pContactUserID int)
BEGIN
	SELECT id FROM pendingcontacts WHERE requestedUserID = pUserID AND addingUserID = pContactUserID;
END $$
DELIMITER ;