-- migrate:up
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE IF NOT EXISTS users (
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  name TEXT NOT NULL,
  email CITEXT NOT NULL UNIQUE,
  address TEXT NOT NULL,
  city TEXT NOT NULL,
  state CITEXT NOT NULL,
  zip TEXT NOT NULL,
  birth_date DATE NOT NULL,
  latitude DOUBLE PRECISION NOT NULL,
  longitude DOUBLE PRECISION NOT NULL,
  password TEXT NOT NULL,
  source CITEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  category CITEXT NOT NULL,
  ean TEXT NOT NULL,
  price DOUBLE PRECISION NOT NULL,
  quantity INT NOT NULL DEFAULT 0,
  rating DOUBLE PRECISION NOT NULL,
  name TEXT NOT NULL,
  vendor TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  user_id BIGINT NOT NULL REFERENCES users(id),
  product_id BIGINT NOT NULL REFERENCES products(id),
  discount DOUBLE PRECISION NOT NULL DEFAULT 0,
  quantity INT NOT NULL,
  subtotal DOUBLE PRECISION NOT NULL,
  tax DOUBLE PRECISION NOT NULL,
  total DOUBLE PRECISION NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
  id BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  reviewer TEXT NOT NULL,
  product_id BIGINT NOT NULL REFERENCES products(id),
  rating INT NOT NULL,
  body TEXT NOT NULL
);

-- migrate:down
