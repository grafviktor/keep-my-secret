package sql

const sqlCreateUserTable = `
CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	login VARCHAR(100) UNIQUE,
	password TEXT NOT NULL,
	restore_password TEXT,
	data_key TEXT
);
`

const sqlCreateSecretTable = `
CREATE TABLE IF NOT EXISTS secret (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	secret_type VARCHAR(10),   -- card, file, pass, note (left non-normalized)
	title TEXT,
	login TEXT,
	password TEXT,
	note TEXT,
	file_name TEXT,
    file BINARY,
	cardholder_name TEXT,
	card_number TEXT,
	expiration TEXT,
	cvv TEXT,
	user_id BIGINT,
	CONSTRAINT fk_secret_user_id FOREIGN KEY(user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
);
`

var sqlInsertUser = `
INSERT INTO user
		(login, password, restore_password, data_key)
	VALUES
		($1, $2, $3, $4)
	RETURNING id;
`

var sqlSelectUser = `
SELECT id, login, password, restore_password, data_key FROM user WHERE login = $1;
`

var sqlInsertSecret = `
INSERT INTO secret (
		secret_type, -- card, file, pass, note (left non-normalized)
		title,
		login,
		password,
		note,
        file,
		file_name,
		cardholder_name,
		card_number,
		expiration,
		cvv,
		user_id
	)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, (SELECT id FROM user WHERE login = $12))
	RETURNING id;
`

var sqlUpdateSecret = `
UPDATE secret SET
		secret_type = $1,
		title = $2,
		login = $3,
		password = $4,
		note = $5,
		file_name = $6,
		cardholder_name = $7,
		card_number = $8,
		expiration = $9,
		cvv = $10
	WHERE id = $11
	AND user_id = (SELECT id FROM user WHERE login = $12);
`

var sqlGetSecretByID = `
SELECT
    id,
    secret_type,
	title,
	login,
	password,
	note,
	file,
	file_name,
	cardholder_name,
	card_number,
	expiration,
	cvv
FROM secret
		WHERE id = $1
		  AND user_id = (
	SELECT id FROM user WHERE login = $2
);
`

var sqlFindSecretsByUser = `
SELECT
    id,
    secret_type,
	title,
	login,
	password,
	note,
	file,
	file_name,
	cardholder_name,
	card_number,
	expiration,
	cvv
FROM secret
	WHERE user_id = (
		SELECT id FROM user WHERE login = $1
);
`

var sqlDeleteSecret = `
DELETE FROM secret
WHERE
    id = $1
AND
    user_id = (
        SELECT id FROM user WHERE login = $2
    );
`
