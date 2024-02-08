
--NOTE: Don't forget to replace AUTOINCREMENT for AUTO_INCREMENT
--      when running this on MySQL instead of Sqlite3
--      (and probably the way the indexes and FKs are declared too)

-- Table: app_config
-- Description: Stores the configuration of the user.
----------------------------------------------------------------------------
-- id -- PK, INTEGER
-- name -- STRING, NOT NULL, UNIQUE, MAXSIZE100
-- value -- STRING, MAXSIZE255
CREATE TABLE IF NOT EXISTS app_config (
'id' INTEGER PRIMARY KEY AUTOINCREMENT,
'name' VARCHAR(100) NOT NULL UNIQUE,
'value' VARCHAR(255)
);
CREATE INDEX IF NOT EXISTS app_config_name ON app_config('name');

-- Table: note
-- Description: stores the user's notes
----------------------------------------------------------------------------
-- id -- PK, INTEGER
-- title -- STRING, NOT NULL, MAXSIZE255
-- content -- STRING, NOT NULL, MAXSIZE131071 (18 bits)
-- content_raw -- STRING, NOT NULL, MAXSIZE131071 (18 bits)
-- views -- INTEGER, NOT NULL, DEFAULT 0
-- lastread_at -- STRING/DATE, MAXSIZE8
-- created_at -- STRING/DATE, NOT NULL, MAXSIZE8
-- updated_at -- STRING/DATE, MAXSIZE8
-- deleted_at -- STRING/DATE, MAXSIZE8
-- deleted -- BOOLEAN, NOT NULL, DEFAULT FALSE
CREATE TABLE IF NOT EXISTS note (
'id' INTEGER PRIMARY KEY AUTOINCREMENT,
'title' VARCHAR(255) NOT NULL,
'content' VARCHAR (131071) NOT NULL,
'content_raw' VARCHAR (131071) NOT NULL,
'views' INTEGER NOT NULL DEFAULT 0,
'lastread_at' VARCHAR(8),
'created_at' VARCHAR(8) NOT NULL,
'updated_at' VARCHAR(8),
'deleted_at' VARCHAR(8),
'deleted' TINYINT(1) NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS note_title ON note('title');
CREATE INDEX IF NOT EXISTS note_content_raw ON note('content_raw');

-- Table: note_dif
-- Description: Stores the user's notes past versions and it's delta values since then.
--              The first entry must be created after the first time the note is edited.
----------------------------------------------------------------------------
-- id -- PK, INTEGER
-- note_id -- FK (note.id), NOT NULL, MAXSIZE255
-- content -- STRING, NOT NULL, MAXSIZE131071 (18 bits)
-- date -- STRING/DATE, NOT NULL, MAXSIZE8
CREATE TABLE IF NOT EXISTS note_dif (
'id' INTEGER PRIMARY KEY AUTOINCREMENT,
'note_id' INTEGER NOT NULL,
'content' VARCHAR(131071),
'date' VARCHAR(8) NOT NULL,
FOREIGN KEY('note_id') REFERENCES note('id')
);
CREATE INDEX IF NOT EXISTS note_dif_note_id ON note_dif('note_id');