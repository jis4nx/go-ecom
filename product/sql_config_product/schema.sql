CREATE TABLE IF NOT EXISTS product(
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  sku VARCHAR(255),
  price NUMERIC(6, 2) NOT NULL,
  available BOOLEAN,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
)
