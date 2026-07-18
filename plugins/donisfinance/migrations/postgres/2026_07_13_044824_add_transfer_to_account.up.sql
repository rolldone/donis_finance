ALTER TABLE transactions ADD COLUMN IF NOT EXISTS to_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_transactions_to_account_id ON transactions(to_account_id);
