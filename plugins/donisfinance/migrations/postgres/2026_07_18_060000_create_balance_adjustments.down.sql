-- +goose Up
CREATE TABLE balance_adjustments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    old_balance BIGINT NOT NULL DEFAULT 0,
    new_balance BIGINT NOT NULL DEFAULT 0,
    reason      TEXT NOT NULL DEFAULT '',
    adjusted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_balance_adjustments_account_id ON balance_adjustments(account_id);

-- +goose Down
DROP TABLE IF EXISTS balance_adjustments;
