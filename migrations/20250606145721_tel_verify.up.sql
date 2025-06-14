SET
    statement_timeout = 0;

--bun:split
CREATE TABLE
    tel_verify (
        tel VARCHAR(10) PRIMARY KEY,
        verify BOOLEAN DEFAULT FALSE
    );

ALTER TABLE users
ADD CONSTRAINT fk_users_tel
FOREIGN KEY (tel)
REFERENCES tel_verify(tel);