CREATE TABLE IF NOT EXISTS containers (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL,
    container_name VARCHAR(255) NOT NULL,
    container_id VARCHAR(255) NOT NULL
);