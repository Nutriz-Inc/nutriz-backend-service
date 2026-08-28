-- DROP TABLES (ordem reversa por FK)
DROP TABLE IF EXISTS route_donation_step;
DROP TABLE IF EXISTS route;
DROP TABLE IF EXISTS bottle;

-- DROP ENUMS
DROP TYPE IF EXISTS enum_route_status;

-- REVERTE enum_user_type (remove 'driver')
-- Falha se existir algum usuario com type = 'driver'.
ALTER TABLE "user" ALTER COLUMN type TYPE VARCHAR(20);
DROP TYPE enum_user_type;
CREATE TYPE enum_user_type AS ENUM ('common', 'adm', 'nurse');
ALTER TABLE "user" ALTER COLUMN type TYPE enum_user_type USING type::enum_user_type;
