DELIMITER $$
CREATE PROCEDURE createAccount(
    pUsername varchar(15),
    pEmailAddress varchar(320),
    pPassword varchar(60),
    pDisplayName varchar(20),
    pValidationGuid varchar(45)
    )
BEGIN
	INSERT INTO users
        (username, password, displayname, email, validationguid)
    VALUES
        (pUsername, pPassword, pDisplayName, pEmailAddress, pValidationGuid);
END $$
DELIMITER ;