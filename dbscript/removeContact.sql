DELIMITER $$
CREATE PROCEDURE removeContact(
    pUserID int,
    pRemoveUserID int
    )
BEGIN

	DELETE FROM friends WHERE userid = pUserID AND frienduserid = pRemoveUserID;
    DELETE FROM friends WHERE userid = pRemoveUserID AND frienduserid = pUserID;
    
END $$
DELIMITER ;