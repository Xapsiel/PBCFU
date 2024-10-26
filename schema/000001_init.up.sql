-- Creating the users table
CREATE TABLE users
(
    id serial NOT NULL UNIQUE,
    login VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    lastclick BIGINT DEFAULT 0,
    permissions int default 0
);

-- Creating the pixels table
CREATE TABLE pixels
(
    x INT,
    y INT,
    id INT,
    color VARCHAR(7) NULL,
    PRIMARY KEY (x, y),
    FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE
);
