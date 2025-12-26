-- This migration is intentionally irreversible.
-- v7 removes invalid SERIAL behavior from foreign keys.
-- Restoring it would reintroduce a schema bug.

RAISE EXCEPTION 'Migration v7 is irreversible';