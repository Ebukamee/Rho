DROP INDEX IF EXISTS products_category_id_idx;

ALTER TABLE products
DROP CONSTRAINT IF EXISTS products_category_id_fkey;

ALTER TABLE products
DROP COLUMN IF EXISTS category_id;