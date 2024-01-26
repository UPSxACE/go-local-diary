
-- Table: app_config
-- Description: stores the configuration of the user
-- id -- PK, INTEGER
-- name -- STRING, NOT NULL, UNIQUE, MAXSIZE100
-- value -- STRING, MAXSIZE255
CREATE TABLE IF NOT EXISTS 'app_config' (
'id' INTEGER PRIMARY KEY,
'name' VARCHAR(100) NOT NULL UNIQUE,
'value' VARCHAR(255)
);