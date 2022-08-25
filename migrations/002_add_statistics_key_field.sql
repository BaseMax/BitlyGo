-- Write your migrate up statements here
alter table links add statistics_key varchar(155) default null;
---- create above / drop below ----
alter table links drop column statistics_key;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
