-- ROUTE: estimated driving/service time, stored as a Go time.Duration (nanoseconds)
ALTER TABLE route ADD COLUMN estimated_time BIGINT;

-- ROUTE DONATION STEP: lifecycle status + free-text description
CREATE TYPE enum_route_donation_step_status AS ENUM ('pending', 'in_progress', 'done', 'error');

ALTER TABLE route_donation_step
  ADD COLUMN status enum_route_donation_step_status NOT NULL DEFAULT 'pending';

ALTER TABLE route_donation_step ADD COLUMN description TEXT;

-- Backfill status for rows that already carry dates or were removed
UPDATE route_donation_step
SET status = (CASE
  WHEN removed_at IS NOT NULL THEN 'error'
  WHEN date_end IS NOT NULL THEN 'done'
  WHEN date_start IS NOT NULL THEN 'in_progress'
  ELSE 'pending'
END)::enum_route_donation_step_status;
