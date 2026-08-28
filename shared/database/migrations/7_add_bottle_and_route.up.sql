-- ENUMS
CREATE TYPE enum_route_status AS ENUM ('pending', 'in_progress', 'done', 'canceled');
ALTER TYPE enum_user_type ADD VALUE IF NOT EXISTS 'driver';

-- BOTTLE
CREATE TABLE bottle (
  id_bottle VARCHAR(36) PRIMARY KEY,

  id_donation VARCHAR(36) NOT NULL,

  quantity_donated_ml NUMERIC(10,2),
  discarded BOOLEAN,
  description TEXT,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,

  CONSTRAINT fk_bottle_donation FOREIGN KEY (id_donation) REFERENCES donation(id_donation),

  CONSTRAINT fk_bottle_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user)
);

-- ROUTE
CREATE TABLE route (
  id_route VARCHAR(36) PRIMARY KEY,

  id_driver VARCHAR(36) NOT NULL,

  name VARCHAR(150) NOT NULL,
  description TEXT NOT NULL,
  user_feedback TEXT,
  city VARCHAR(100),
  neighborhood VARCHAR(100),
  status enum_route_status NOT NULL,
  date_start TIMESTAMP,
  date_end TIMESTAMP,
  mileage NUMERIC(10,2),
  date_set TIMESTAMP NOT NULL,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  removed_at TIMESTAMP,
  removed_by VARCHAR(36),

  CONSTRAINT fk_route_driver FOREIGN KEY (id_driver) REFERENCES "user"(id_user),

  CONSTRAINT fk_route_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_route_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_route_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

-- ROUTE DONATION STEP
CREATE TABLE route_donation_step (
  id_route_donation_step VARCHAR(36) PRIMARY KEY,

  id_route VARCHAR(36) NOT NULL,
  id_donation_step VARCHAR(36) NOT NULL,

  stop_order INTEGER,
  date_start TIMESTAMP,
  date_end TIMESTAMP,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  removed_at TIMESTAMP,
  removed_by VARCHAR(36),

  CONSTRAINT fk_rds_route FOREIGN KEY (id_route) REFERENCES route(id_route),
  CONSTRAINT fk_rds_donation_step FOREIGN KEY (id_donation_step) REFERENCES donation_step(id_donation_step),

  CONSTRAINT fk_rds_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_rds_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_rds_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);
