-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  total_sum DECIMAL(19, 4) NOT NULL,
  status_code SMALLINT NOT NULL DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  confirmed_at TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT chk_status_code CHECK (status_code IN (0, 1, 2, 3))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments
-- +goose StatementEnd
