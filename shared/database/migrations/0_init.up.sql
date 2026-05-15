-- ENUMS
CREATE TYPE enum_user_type AS ENUM ('common', 'adm', 'nurse');
CREATE TYPE enum_donation_step_status AS ENUM ('pending', 'review', 'done', 'warn', 'failed');
CREATE TYPE enum_donation_steps as ENUM (
  'Exame de sangue',
  'Entregar kit de ordenha',
  'Coletar leite',
  'Análise de leite'
);
CREATE TYPE enum_job_status AS ENUM ('pending', 'done', 'failed');

-- USER
CREATE TABLE "user" (
  id_user VARCHAR(36) PRIMARY KEY,

  internal_identifier VARCHAR(36) UNIQUE,
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
  neighborhood VARCHAR(100) NOT NULL,
  complement VARCHAR(100),
  latitude NUMERIC(10,7),
  longitude NUMERIC(10,7),

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  removed_at TIMESTAMP,

  CONSTRAINT fk_address_user FOREIGN KEY (id_user) REFERENCES "user"(id_user),
  CONSTRAINT fk_address_dp FOREIGN KEY (id_donation_point) REFERENCES donation_point(id_donation_point)
);

-- DONATION
CREATE TABLE donation (
  id_donation VARCHAR(36) PRIMARY KEY,

  quantity_donated NUMERIC(10,2),
  is_active BOOLEAN NOT NULL,
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
  id_donation_step VARCHAR(36) PRIMARY KEY,

  id_donation VARCHAR(36) NOT NULL,

  name enum_donation_steps NOT NULL,
  description TEXT NOT NULL,
  status enum_donation_step_status NOT NULL,
  set_date TIMESTAMP,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,
  updated_at TIMESTAMP,
  updated_by VARCHAR(36),
  completed_at TIMESTAMP,

  CONSTRAINT fk_step_donation FOREIGN KEY (id_donation) REFERENCES donation(id_donation),

  CONSTRAINT fk_donation_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_donation_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user)
);

-- DONATION STEP TIMELINE
CREATE TABLE donation_step_timeline (
  id_donation_step_timeline VARCHAR(36) PRIMARY KEY,

  id_donation_step VARCHAR(36) NOT NULL,

  description TEXT NOT NULL,
  status enum_donation_step_status NOT NULL,
  set_date TIMESTAMP,

  created_at TIMESTAMP NOT NULL,
  created_by VARCHAR(36) NOT NULL,

  CONSTRAINT fk_timeline_donation_step FOREIGN KEY (id_donation_step) REFERENCES donation_step(id_donation_step),

  CONSTRAINT fk_donation_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user)
);

-- JOB
CREATE TABLE job (
  id_job VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36) NOT NULL,
  id_step VARCHAR(36) NOT NULL,

  status enum_job_status NOT NULL,
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
  CONSTRAINT fk_job_step FOREIGN KEY (id_step) REFERENCES donation_step(id_donation_step),

  CONSTRAINT fk_job_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

-- CONSENT LOG
CREATE TABLE consent_log (
  id_consent_log VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36) NOT NULL,

  terms_version VARCHAR(50) NOT NULL,
  accepted_at TIMESTAMP NOT NULL,
  ip_address VARCHAR(45) NOT NULL,
  user_agent TEXT NOT NULL,

  CONSTRAINT fk_consent_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);

-- USER BABY
CREATE TABLE user_baby (
  id_user_baby VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36) NOT NULL,

  name VARCHAR(120),
  birth_date TIMESTAMP NOT NULL,

  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  removed_at TIMESTAMP,

  CONSTRAINT fk_user_baby_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);