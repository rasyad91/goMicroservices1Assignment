USE assignment5;
CREATE TABLE courses (
    id INT PRIMARY KEY,
    title VARCHAR(30),
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);