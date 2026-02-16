CREATE TABLE IF NOT EXISTS accounts (
    account_id SERIAL PRIMARY KEY,
    document_number TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS operations_types (
    operation_type_id SMALLINT PRIMARY KEY,
    description TEXT NOT NULL
);

INSERT INTO operations_types (operation_type_id, description) VALUES
    (1, 'PURCHASE'),
    (2, 'INSTALLMENT PURCHASE'),
    (3, 'WITHDRAWAL'),
    (4, 'PAYMENT')
ON CONFLICT (operation_type_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS transactions (
    transaction_id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL REFERENCES accounts(account_id),
    operation_type_id SMALLINT NOT NULL REFERENCES operations_types(operation_type_id),
    amount NUMERIC(12,2) NOT NULL,
    event_date TIMESTAMP WITH TIME ZONE NOT NULL
)
