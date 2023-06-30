CREATE TABLE IF NOT EXISTS Users (
	ID      	INTEGER    NOT NULL PRIMARY KEY,
	Created 	INTEGER 	NOT NULL DEFAULT (strftime('%s', 'now')) -- unix seconds
);