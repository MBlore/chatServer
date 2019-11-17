DELIMITER $$
CREATE PROCEDURE removeContact(
    pUserID int,
    pRemoveUserID int
    )
BEGIN

	DELETE FROM friends WHERE userid = pUserID AND frienduserid = pRemoveUserID;

END $$
DELIMITER ;