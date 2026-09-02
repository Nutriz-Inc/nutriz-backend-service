ALTER TABLE route_donation_step DROP COLUMN IF EXISTS description;
ALTER TABLE route_donation_step DROP COLUMN IF EXISTS status;
DROP TYPE IF EXISTS enum_route_donation_step_status;
ALTER TABLE route DROP COLUMN IF EXISTS estimated_time;
