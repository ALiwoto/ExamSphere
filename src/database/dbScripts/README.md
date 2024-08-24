# ExamSphere SQL Database Scripts

This folder contains the SQL scripts to create the database schema and populate the tables with data.

## WARNINGS
  - **The scripts are written for PostgreSQL. They may not work for other databases.**
  - **Do NOT run these migrations directly yourself. The application should run them automatically if necessary.**


## Run the migrations
Run migrations from 1 to 3 in the given order to create the database schema.

```bash
psql -U postgres --dbname examspheredb -a -f migration1.sql
psql -U postgres --dbname examspheredb -a -f migration2.sql
psql -U postgres --dbname examspheredb -a -f migration3.sql
```

## Drop the defined types and functions
Run the following commands to drop the defined types and functions.

```bash
psql -U postgres --dbname examspheredb -a -f drop_all.sql
```

## SQL Features used in the scripts
  - [x] DDL (Data Definition Language) commands like CREATE, ALTER, DROP, TRUNCATE
  - [x] DML (Data Manipulation Language) commands like INSERT, UPDATE, DELETE
  - [ ] DCL (Data Control Language) commands like GRANT, REVOKE
  - [x] TCL (Transaction Control Language) commands like COMMIT, ROLLBACK, SAVEPOINT

## Differences between FUNCTION and PROCEDURE
  - [x] FUNCTION returns a value, PROCEDURE does not.
  - [x] FUNCTION can be used in SELECT statement, PROCEDURE cannot.
  - [x] FUNCTIONs can be called from procedures, but procedures cannot be called from functions. 
  - [x] In PostgreSQL, FUNCTIONS should not have TRANSACTION control commands like COMMIT, ROLLBACK, SAVEPOINT, etc. PROCEDUREs can have these commands; because FUNCTIONs themselves act like transactions.


## SQL Statements used in the scripts
  - [x] CREATE DATABASE
  - [x] CREATE TABLE
  - [x] ALTER TABLE
  - [x] DROP TABLE
  - [ ] TRUNCATE TABLE
  - [x] INSERT INTO
  - [x] UPDATE
  - [x] DELETE
  - [ ] GRANT
  - [ ] REVOKE
  - [x] COMMIT
  - [x] ROLLBACK
  - [ ] SAVEPOINT
  - [x] SELECT
  - [x] WHERE
  - [x] ORDER BY
  - [x] GROUP BY
  - [x] JOIN
  - [ ] UNION
  - [ ] INTERSECT
  - [ ] EXCEPT
  - [ ] HAVING
  - [x] LIMIT
  - [x] OFFSET
  - [x] LIKE, ILIKE
  - [x] IN
  - [ ] BETWEEN
  - [x] EXISTS
  - [x] NOT EXISTS
  - [ ] ANY
  - [ ] ALL
  - [x] WHEN, CASE
  - [ ] COALESCE
  - [ ] CAST
  - [x] COMMENT
  - [ ] IS NULL, IS NOT NULL, NULLIF, NOT NULL
  - [ ] IS DISTINCT FROM, IS NOT DISTINCT FROM
  - [x] VIEW, CREATE VIEW, DROP VIEW
  - [ ] INDEX, CREATE INDEX, DROP INDEX
  - [ ] SEQUENCE, CREATE SEQUENCE, DROP SEQUENCE
  - [x] TRIGGER, CREATE TRIGGER, DROP TRIGGER
  - [ ] STORED PROCEDURE, CREATE PROCEDURE, DROP PROCEDURE
  - [ ] STORED FUNCTION, CREATE FUNCTION, DROP FUNCTION
  - [ ] STORED TRIGGER, CREATE TRIGGER, DROP TRIGGER
  - [ ] STORED PACKAGE, CREATE PACKAGE, DROP PACKAGE

## SQL Data Types used in the scripts
  - [x] INTEGER
  - [x] BIGINT
  - [x] SMALLINT
  - [x] DECIMAL
  - [x] NUMERIC
  - [x] REAL
  - [x] DOUBLE PRECISION
  - [ ] CHAR
  - [x] VARCHAR
  - [x] TEXT
  - [x] DATE
  - [x] TIME
  - [x] TIMESTAMP
  - [x] BOOLEAN
  - [ ] BINARY
  - [ ] VARBINARY
  - [ ] BLOB
  - [ ] CLOB
  - [ ] ARRAY
  - [ ] JSON
  - [ ] XML
  - [ ] UUID
  - [ ] ENUM
  - [ ] GEOMETRY
  - [ ] POINT
  - [ ] LINESTRING
  - [ ] POLYGON
  - [ ] MULTIPOINT
  - [ ] MULTILINESTRING
  - [ ] MULTIPOLYGON
  - [ ] GEOMETRYCOLLECTION
  - [ ] INT2
  - [ ] INT4
  - [ ] INT8
  - [x] SERIAL
  - [ ] BIGSERIAL
  - [ ] MONEY
  - [ ] OID
  - [ ] REGCLASS
  - [ ] REGCONFIG

## SQL Constraints used in the scripts
  - [x] PRIMARY KEY
  - [x] FOREIGN KEY
  - [x] UNIQUE
  - [x] CHECK
  - [x] NOT NULL
  - [x] DEFAULT
  - [ ] INDEX
  - [ ] EXCLUSION
  - [ ] PARTIAL
  - [ ] DEFERRABLE
  - [ ] INITIALLY DEFERRED
  - [ ] INITIALLY IMMEDIATE
  - [ ] MATCH
  - [x] ON DELETE
  - [x] ON UPDATE
  - [x] REFERENCES
  - [x] TRIGGER
  - [x] VIEW
  - [x] STORED PROCEDURE
  - [x] FUNCTION

## Resources used (Useful links)
  - [PostgreSQL Documentation](https://www.postgresql.org/docs/)
  - [SQL Tutorial](https://www.w3schools.com/sql/)
  - [SQL Cheat-sheet](https://learnsql.com/blog/sql-cheat-sheet/)
  - [SQL Fiddle](http://sqlfiddle.com/)
  - [SQL Formatter](https://sqlformat.org/)
  - [ON CONFLICT DO NOTHING](https://www.prisma.io/dataguide/postgresql/inserting-and-modifying-data/insert-on-conflict)
  - [Stored procedure Vs. Function: What are the differences?](https://www.shiksha.com/online-courses/articles/stored-procedure-vs-function-what-are-the-differences/)
  - [Difference between Stored Procedure and Function in SQL Server](https://www.scholarhat.com/tutorial/sqlserver/difference-between-stored-procedure-and-function-in-sql-server)
  - [Can we call a procedure in select statement with any restriction?](https://asktom.oracle.com/ords/asktom.search?tag=can-we-call-a-procedure-in-select-statement-with-any-restriction&p_session=606595063924099#:~:text=The%20execution%20of%20a%20function,it%20does%20not%20return%20anything.)