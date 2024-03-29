DELIMITER $$
CREATE PROCEDURE confirmContact(
    pRequestingUserID int,
    pAddingUserID int
    )
BEGIN

	DELETE FROM pendingcontacts WHERE requestedUserID = pRequestingUserID AND addingUserID = pAddingUserID;

    INSERT INTO friends (userid, frienduserid) VALUES (pRequestingUserID, pAddingUserID);
    INSERT INTO friends (frienduserid, userid) VALUES (pRequestingUserID, pAddingUserID);

END $$
DELIMITER ;