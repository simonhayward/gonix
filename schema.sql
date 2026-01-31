CREATE TABLE IF NOT EXISTS User (
    UserID   INTEGER NOT NULL AS (Data->>'user_id')  STORED UNIQUE, -- PK
    Created  INTEGER NOT NULL AS (unixepoch(Data->>'created')) STORED,
    Data     JSONB   NOT NULL
);

CREATE TABLE IF NOT EXISTS Token (
  Token    TEXT    NOT NULL AS (Data->>'token')  STORED UNIQUE, -- PK
  UserID   INTEGER NOT NULL AS (Data->>'user_id') STORED REFERENCES User (UserID),
  Created  INTEGER NOT NULL AS (unixepoch(Data->>'created')) STORED,
  LastUsed INTEGER AS (unixepoch(Data->>'last_used')) CHECK (LastUsed>0),
  Data     JSONB   NOT NULL
);
