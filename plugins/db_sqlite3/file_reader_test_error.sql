-- this line shall be ignored
-- this line shall be ignored
CREATE TABLE 'table1'(
'id' INTEGER PRIMARY KEY,
'name' VARCHAR(12)
);

INSERT INTO 'table1'('name') VALUES ('John');
INSERT INTO 'table1'('name') VALUES ('Maria')
INSERT INTO 'table1'('name') VALUES ('Curry');
-- should end up in error because of lack of semicolon on line 9
-- these 2 lines cannot be reached so they should not count as ignored