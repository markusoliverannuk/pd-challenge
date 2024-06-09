CREATE TABLE
    IF NOT EXISTS gists (
        username TEXT PRIMARY KEY UNIQUE NOT NULL,
        description TEXT
    );

CREATE TABLE
    IF NOT EXISTS files (
        username TEXT NOT NULL,
        path TEXT NOT NULL,
        PRIMARY KEY (username, path),
        FOREIGN KEY (username) REFERENCES gists (username)
    )