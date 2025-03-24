CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS versions (
    id SERIAL PRIMARY KEY,
    service_id INTEGER NOT NULL,
    number TEXT,
    FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE
);