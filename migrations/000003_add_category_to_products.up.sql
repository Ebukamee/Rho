ALTER TABLE products
ADD COLUMN category_id UUID;

ALTER TABLE products
ADD CONSTRAINT products_category_id_fkey
FOREIGN KEY (category_id)
REFERENCES categories(id)
ON DELETE SET NULL;

CREATE INDEX products_category_id_idx
ON products(category_id);