-- DROP TABLES (ordem reversa por FK)

DROP TABLE IF EXISTS file;
DROP TABLE IF EXISTS job;
DROP TABLE IF EXISTS donation_step;
DROP TABLE IF EXISTS donation;
DROP TABLE IF EXISTS address;
DROP TABLE IF EXISTS donation_point;
DROP TABLE IF EXISTS "user";

-- DROP ENUMS
DROP TYPE IF EXISTS enum_donation_step_status;
DROP TYPE IF EXISTS enum_user_type;