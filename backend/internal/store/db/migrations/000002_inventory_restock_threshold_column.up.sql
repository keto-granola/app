ALTER TABLE inventory
ADD COLUMN IF NOT EXISTS restock_threshold INTEGER NOT NULL CHECK (restock_threshold > 0);