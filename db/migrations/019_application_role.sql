-- Write your migrate up statements here
CREATE ROLE socio_app_user WITH LOGIN PASSWORD '38d03145-5d69-4690-9f1d-293daa874927';
GRANT ALL PRIVILEGES ON DATABASE socio TO socio_app_user;
---- create above / drop below ----
REVOKE ALL PRIVILEGES ON DATABASE socio FROM socio_app_user;
DROP ROLE socio_app_user;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
