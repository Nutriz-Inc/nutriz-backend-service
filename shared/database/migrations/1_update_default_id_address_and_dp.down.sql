ALTER TABLE donation_point ALTER COLUMN id_donation_point DROP DEFAULT;
ALTER TABLE address ALTER COLUMN id_address DROP DEFAULT;

DROP FUNCTION IF EXISTS generate_ksuid();