ALTER TABLE "donation_step_timeline"
DROP CONSTRAINT IF EXISTS fk_donation_step_timeline_id_address,
DROP COLUMN IF EXISTS id_address;

ALTER TABLE "donation_step"
DROP CONSTRAINT IF EXISTS fk_donation_step_id_address,
DROP COLUMN IF EXISTS id_address;