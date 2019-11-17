DELIMITER $$
CREATE PROCEDURE addPendingContact(
    pRequestingUserID int,
    pAddingUserID int,
    pMessage varchar(100)
    )
BEGIN
	INSERT INTO pendingcontacts
        (requestedUserID, addingUserID, message)
    VALUES
        (pRequestingUserID, pAddingUserID, pMessage);
END $$
DELIMITER ;