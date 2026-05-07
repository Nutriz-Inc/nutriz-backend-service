-- ENUMS
CREATE TYPE enum_user_type AS ENUM ('common', 'adm', 'nurse');
CREATE TYPE enum_donation_step_status AS ENUM ('pending', 'review', 'done', 'warn', 'failed');

-- USER
CREATE TABLE "user" (
  id_user VARCHAR(36) PRIMARY KEY,

  type enum_user_type NOT NULL,

  name VARCHAR(120) NOT NULL,
  cpf VARCHAR(11) NOT NULL UNIQUE,
  birth_date TIMESTAMP NOT NULL,
  phone_number VARCHAR(20) NOT NULL UNIQUE,

  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  milk_donated NUMERIC,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  removed_at TIMESTAMP,
  removed_by VARCHAR(36),

  CONSTRAINT fk_user_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_user_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

-- DONATION POINT
CREATE TABLE donation_point (
  id_donation_point VARCHAR(36) PRIMARY KEY,

  name VARCHAR(150) NOT NULL,
  description TEXT,

  has_home BOOLEAN NOT NULL,
  phone_number VARCHAR(20),
  email VARCHAR(255),
  opening_hours VARCHAR(255),

  removed_at TIMESTAMP
);

-- ADDRESS
CREATE TABLE address (
  id_address VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36),
  id_donation_point VARCHAR(36),

  zipcode VARCHAR(10) NOT NULL, -- CEP
  street VARCHAR(150) NOT NULL,
  number VARCHAR(10),
  city VARCHAR(100) NOT NULL,
  state VARCHAR(2) NOT NULL,
  complement VARCHAR(100),
  latitude NUMERIC(10,7),
  longitude NUMERIC(10,7),

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  removed_at TIMESTAMP,
  removed_by VARCHAR(36),

  CONSTRAINT fk_address_user FOREIGN KEY (id_user) REFERENCES "user"(id_user),
  CONSTRAINT fk_address_dp FOREIGN KEY (id_donation_point) REFERENCES donation_point(id_donation_point),

  CONSTRAINT fk_address_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_address_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_address_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

-- DONATION
CREATE TABLE donation (
  id_donation VARCHAR(36) PRIMARY KEY,

  is_active BOOLEAN NOT NULL,

  quantity NUMERIC(10,2),
  user_feedback TEXT,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  removed_at TIMESTAMP,
  removed_by VARCHAR(36),

  CONSTRAINT fk_donation_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_donation_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_donation_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

-- DONATION STEP
CREATE TABLE donation_step (
  id_step_donation VARCHAR(36) PRIMARY KEY,

  id_donation VARCHAR(36) NOT NULL,

  name VARCHAR(100) NOT NULL,
  description TEXT NOT NULL,

  status enum_donation_step_status NOT NULL,

  set_date TIMESTAMP,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  completed_at TIMESTAMP,

  CONSTRAINT fk_step_donation FOREIGN KEY (id_donation) REFERENCES donation(id_donation)
);

-- JOB
CREATE TABLE job (
  id_job VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36),
  id_step VARCHAR(36),

  name VARCHAR(120) NOT NULL,
  description TEXT NOT NULL,
  date_set TIMESTAMP,

  user_feedback TEXT,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  removed_at TIMESTAMP,
  removed_by VARCHAR(36),

  CONSTRAINT fk_job_user FOREIGN KEY (id_user) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_step FOREIGN KEY (id_step) REFERENCES donation_step(id_step_donation),

  CONSTRAINT fk_job_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

-- FILE
CREATE TABLE file (
  id_file VARCHAR(36) PRIMARY KEY,

  id_job VARCHAR(36),

  file_path VARCHAR(500) NOT NULL,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  removed_at TIMESTAMP,
  removed_by VARCHAR(36),

  CONSTRAINT fk_file_job FOREIGN KEY (id_job) REFERENCES job(id_job),

  CONSTRAINT fk_file_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_file_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_file_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);