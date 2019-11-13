DELIMITER $$
CREATE PROCEDURE getFriends(pUserID int)
BEGIN
	SELECT
        f.frienduserid,
        u.username,
        u.displayname,
        u.status,
        u.userimage,
        u.statustext
    FROM
        chatzorz.friends AS f
        JOIN chatzorz.users AS u ON u.id = f.frienduserid
    WHERE
        f.userid = pUserID;
END $$
DELIMITER ;