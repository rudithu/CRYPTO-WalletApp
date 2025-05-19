INSERT INTO ccy_conversion (from_ccy, to_ccy, rate)
VALUES 
  ('USD', 'SGD', 1.35),
  ('USD', 'EUR', 0.92),
  ('USD', 'JPY', 155.25),
  ('USD', 'GBP', 0.79),
  ('USD', 'AUD', 1.52),
  ('USD', 'CHF', 0.91),
  ('USD', 'CAD', 1.36),
  ('USD', 'BTC', 0.000010);

INSERT INTO users (id, name)
VALUES
    (1, 'Alice'),
    (2, 'Bob'),
    (3, 'Charlie'),
    (4, 'Danny');

INSERT INTO wallets (user_id, balance, currency, type, is_default)
VALUES
    (1, 0, 'USD', 'saving', true),
    (2, 0, 'SGD', 'trading', false),
    (2, 0, 'AUD', 'saving', true),
    (2, 0, 'USD', 'saving', false),
    (3, 0, 'BTC', 'trading', false),
    (3, 0, 'JPY', 'saving', false),
    (3, 0, 'USD', 'saving', true),
    (4, 0, 'USD', 'saving', false),
    (4, 0, 'SGD', 'trading', true);
    