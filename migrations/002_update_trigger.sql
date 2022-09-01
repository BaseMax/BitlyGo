-- Write your migrate up statements here
CREATE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    New.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
---- create above / drop below ----
DROP FUNCTION update_updated_at_column;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
