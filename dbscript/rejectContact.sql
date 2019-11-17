DELIMITER $$
CREATE PROCEDURE rejectContact(
    pRequestingUserID int,
    pAddingUserID int
    )
BEGIN

	DELETE FROM pendingcontacts WHERE requestedUserID = pRequestingUserID AND addingUserID = pAddingUserID;

END $$
DELIMITER ;