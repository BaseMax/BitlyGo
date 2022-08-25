-- Write your migrate up statements here
alter table links alter column owner_id drop not null;
---- create above / drop below ----
alter table links alter column owner_id set not null;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
