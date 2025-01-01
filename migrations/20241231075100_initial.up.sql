CREATE TABLE IF NOT EXISTS authors (
    id serial PRIMARY KEY,
    name varchar(65) NOT NULL,
    email varchar(128) NOT NULL UNIQUE,
    birth_date date NULL,
    bio TEXT NULL
);

CREATE TABLE IF NOT EXISTS books (
    id serial PRIMARY KEY,
    title varchar(128) NOT NULL,
    description TEXT NULL,
    publish_date date NOT NULL,
    author_id serial NOT NULL,
    CONSTRAINT author_books FOREIGN KEY(author_id) REFERENCES authors(id) ON DELETE CASCADE
);
