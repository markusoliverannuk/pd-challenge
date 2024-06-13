CREATE TABLE
    IF NOT EXISTS gists (
        id TEXT UNIQUE NOT NULL,
        username TEXT NOT NULL,
        description TEXT,
        seen INTEGER
    );

CREATE TABLE
    IF NOT EXISTS files (
        id TEXT NOT NULL,
        username TEXT NOT NULL,
        path TEXT NOT NULL,
        PRIMARY KEY (username, path),
        FOREIGN KEY (username) REFERENCES gists (username)
    )