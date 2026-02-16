CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    document_number TEXT UNIQUE NOT NULL
);
