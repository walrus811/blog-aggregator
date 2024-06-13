-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT now() NOT NULL,
  name VARCHAR  NOT NULL,
  url VARCHAR UNIQUE NOT NULL,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
Drop TABLE feeds;