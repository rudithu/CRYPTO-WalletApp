CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    currency TEXT NOT NULL,         -- e.g., BTC, ETH, USD
    type TEXT,                      -- e.g., saving, trading, cold-storage
    balance NUMERIC(20, 2) DEFAULT 0.00,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Composite unique constraint
    CONSTRAINT unique_user_currency UNIQUE (user_id, currency)
);

-- Partial unique index: only one default wallet per user
CREATE UNIQUE INDEX IF NOT EXISTS one_default_wallet_per_user
ON wallets(user_id)
WHERE is_default = TRUE;

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    wallet_id INT REFERENCES wallets(id),
    type VARCHAR(20), -- deposit, withdrawal, transfer
    amount NUMERIC(20, 2),
    counterparty_wallet_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ccy_conversion (
    from_ccy TEXT NOT NULL,
    to_ccy TEXT NOT NULL,
    rate NUMERIC(20, 6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (from_ccy, to_ccy)
);

INSERT INTO ccy_conversion (from_ccy, to_ccy, rate)
VALUES 
  ('USD', 'SGD', 1.35),
  ('USD', 'EUR', 0.92),
  ('USD', 'JPY', 155.25),
  ('USD', 'GBP', 0.79),
  ('USD', 'AUD', 1.52),
  ('USD', 'CHF', 0.91),
  ('USD', 'CAD', 1.36);