DELIMITER $$
CREATE PROCEDURE getUserPendingContacts(pUserId int)
BEGIN
    SELECT
        u.id,
        u.username,
        u.displayname,
        u.userimage,
        pc.message
    FROM
        chatzorz.pendingcontacts AS pc
        JOIN chatzorz.users AS u ON u.id = pc.requestedUserID
    WHERE
        pc.addingUserID = pUserID;

END $$
DELIMITER ;