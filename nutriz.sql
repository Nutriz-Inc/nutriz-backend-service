--Criação do banco de dados (DDL):

--------------------------------------------------------------------------------
-- EXTENSÕES
--------------------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS pgcrypto; -- gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS vector;   -- pgvector (tipo VECTOR e índice HNSW)

--------------------------------------------------------------------------------
-- ENUMS
--------------------------------------------------------------------------------
CREATE TYPE enum_user_type AS ENUM ('common', 'adm', 'nurse');
CREATE TYPE enum_donation_step_status AS ENUM ('pending', 'review', 'done', 'warn', 'failed');
CREATE TYPE enum_donation_steps AS ENUM (
  'Exame de sangue',
  'Entregar kit de ordenha',
  'Coletar leite',
  'Análise de leite'
);
CREATE TYPE enum_job_status AS ENUM ('pending', 'done', 'failed');

--------------------------------------------------------------------------------
-- FUNÇÃO AUXILIAR: gerador de KSUID (base62, 27 caracteres, time-sortable)
-- usado como DEFAULT para donation_point.id_donation_point e address.id_address.
--------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION generate_ksuid() RETURNS text AS $$
DECLARE
    v_time timestamp with time zone := NULL;
    v_seconds numeric(50) := NULL;
    v_numeric numeric(50) := NULL;
    v_epoch numeric(50) := 1400000000; -- 2014-05-13T16:53:20Z
    v_base62 text := '';
    v_alphabet text[] := ARRAY[
        '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
        'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
        'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
        'U', 'V', 'W', 'X', 'Y', 'Z',
        'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
        'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
        'u', 'v', 'w', 'x', 'y', 'z'];
BEGIN
    v_time := clock_timestamp();
    v_seconds := floor(EXTRACT(EPOCH FROM v_time)) - v_epoch;
    v_numeric := v_seconds * pow(2::numeric(50), 128)
        + ((random()::numeric(70,20) * pow(2::numeric(70,20), 48))::numeric(50) * pow(2::numeric(50), 80)::numeric(50))
        + ((random()::numeric(70,20) * pow(2::numeric(70,20), 40))::numeric(50) * pow(2::numeric(50), 40)::numeric(50))
        +  (random()::numeric(70,20) * pow(2::numeric(70,20), 40))::numeric(50);

    WHILE v_numeric <> 0 LOOP
        v_base62 := v_base62 || v_alphabet[mod(v_numeric, 62)::int + 1];
        v_numeric := div(v_numeric, 62);
    END LOOP;
    v_base62 := reverse(v_base62);
    v_base62 := lpad(v_base62, 27, '0');

    RETURN v_base62;
END $$ LANGUAGE plpgsql;

--------------------------------------------------------------------------------
-- USER
--------------------------------------------------------------------------------
CREATE TABLE "user" (
  id_user              VARCHAR(36)  NOT NULL,
  internal_identifier  VARCHAR(36),
  type                 enum_user_type NOT NULL,
  name                 VARCHAR(120) NOT NULL,
  cpf                  VARCHAR(11)  NOT NULL,
  birth_date           TIMESTAMP    NOT NULL,
  phone_number         VARCHAR(20)  NOT NULL,
  email                VARCHAR(255) NOT NULL,
  password             VARCHAR(255) NOT NULL,
  milk_donated         NUMERIC,

  created_at           TIMESTAMP    NOT NULL,
  created_by           VARCHAR(36)  NOT NULL,
  updated_at           TIMESTAMP,
  updated_by           VARCHAR(36),
  removed_at           TIMESTAMP,
  removed_by           VARCHAR(36),

  CONSTRAINT pk_user PRIMARY KEY (id_user),
  CONSTRAINT fk_user_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_user_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
  -- migration 6 removeu as UNIQUE de internal_identifier, cpf, phone_number
  -- e email desta tabela; propositalmente não recriadas aqui.
);

--------------------------------------------------------------------------------
-- DONATION POINT
--------------------------------------------------------------------------------
CREATE TABLE donation_point (
  id_donation_point VARCHAR(36)  NOT NULL DEFAULT ('dpt_' || generate_ksuid()),

  name               VARCHAR(150) NOT NULL,
  description        TEXT,
  has_home           BOOLEAN      NOT NULL,
  phone_number       VARCHAR(20),
  email              VARCHAR(255),
  opening_hours      VARCHAR(255),

  removed_at         TIMESTAMP,

  CONSTRAINT pk_donation_point PRIMARY KEY (id_donation_point),
  CONSTRAINT uq_donation_point_email UNIQUE (email)
);

--------------------------------------------------------------------------------
-- ADDRESS
--------------------------------------------------------------------------------
CREATE TABLE address (
  id_address        VARCHAR(36)  NOT NULL DEFAULT ('adr_' || generate_ksuid()),

  id_user           VARCHAR(36),
  id_donation_point VARCHAR(36),

  zipcode           VARCHAR(10)  NOT NULL, -- CEP
  street            VARCHAR(150) NOT NULL,
  number            VARCHAR(10),
  city              VARCHAR(100) NOT NULL,
  state             VARCHAR(2)   NOT NULL,
  neighborhood      VARCHAR(100) NOT NULL,
  complement        VARCHAR(100),
  latitude          NUMERIC(10,7),
  longitude         NUMERIC(10,7),

  created_at        TIMESTAMP    NOT NULL,
  updated_at        TIMESTAMP,
  removed_at        TIMESTAMP,

  CONSTRAINT pk_address PRIMARY KEY (id_address),
  CONSTRAINT ck_address_state CHECK (LENGTH(state) = 2),
  CONSTRAINT fk_address_user FOREIGN KEY (id_user) REFERENCES "user"(id_user),
  CONSTRAINT fk_address_dp FOREIGN KEY (id_donation_point) REFERENCES donation_point(id_donation_point)
);

--------------------------------------------------------------------------------
-- DONATION
--------------------------------------------------------------------------------
CREATE TABLE donation (
  id_donation       VARCHAR(36) NOT NULL,

  quantity_donated  NUMERIC(10,2),
  is_active         BOOLEAN     NOT NULL,
  user_feedback     TEXT,
  score_feedback    SMALLINT,                          -- migration 5

  created_at        TIMESTAMP   NOT NULL,
  created_by        VARCHAR(36) NOT NULL,
  updated_at        TIMESTAMP,
  updated_by        VARCHAR(36),
  removed_at        TIMESTAMP,
  removed_by        VARCHAR(36),

  CONSTRAINT pk_donation PRIMARY KEY (id_donation),
  CONSTRAINT ck_donation_quantity CHECK (quantity_donated >= 0),
  CONSTRAINT ck_donation_score_feedback CHECK (score_feedback BETWEEN 1 AND 5),
  CONSTRAINT fk_donation_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_donation_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_donation_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- DONATION STEP
--------------------------------------------------------------------------------
CREATE TABLE donation_step (
  id_donation_step VARCHAR(36) NOT NULL,

  id_donation      VARCHAR(36) NOT NULL,
  id_address       VARCHAR(36),                         -- migration 4

  name             enum_donation_steps NOT NULL,
  description      TEXT         NOT NULL,
  status           enum_donation_step_status NOT NULL,
  set_date         TIMESTAMP,

  created_at       TIMESTAMP   NOT NULL,
  created_by       VARCHAR(36) NOT NULL,
  updated_at       TIMESTAMP,
  updated_by       VARCHAR(36),
  completed_at     TIMESTAMP,

  CONSTRAINT pk_donation_step PRIMARY KEY (id_donation_step),
  CONSTRAINT fk_donation_step_donation FOREIGN KEY (id_donation) REFERENCES donation(id_donation),
  CONSTRAINT fk_donation_step_address FOREIGN KEY (id_address) REFERENCES address(id_address),
  CONSTRAINT fk_donation_step_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_donation_step_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- DONATION STEP TIMELINE
--------------------------------------------------------------------------------
CREATE TABLE donation_step_timeline (
  id_donation_step_timeline VARCHAR(36) NOT NULL,

  id_donation_step          VARCHAR(36) NOT NULL,
  id_address                VARCHAR(36),                -- migration 4

  description               TEXT        NOT NULL,
  status                    enum_donation_step_status NOT NULL,
  set_date                  TIMESTAMP,

  created_at                TIMESTAMP   NOT NULL,
  created_by                VARCHAR(36) NOT NULL,

  CONSTRAINT pk_donation_step_timeline PRIMARY KEY (id_donation_step_timeline),
  CONSTRAINT fk_timeline_donation_step FOREIGN KEY (id_donation_step) REFERENCES donation_step(id_donation_step),
  CONSTRAINT fk_timeline_address FOREIGN KEY (id_address) REFERENCES address(id_address),
  CONSTRAINT fk_timeline_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- JOB
--------------------------------------------------------------------------------
CREATE TABLE job (
  id_job         VARCHAR(36) NOT NULL,

  id_user        VARCHAR(36) NOT NULL,
  id_step        VARCHAR(36) NOT NULL,

  status         enum_job_status NOT NULL,
  name           VARCHAR(120) NOT NULL,
  description    TEXT         NOT NULL,
  date_set       TIMESTAMP,
  user_feedback  TEXT,

  created_at     TIMESTAMP   NOT NULL,
  created_by     VARCHAR(36) NOT NULL,
  updated_at     TIMESTAMP,
  updated_by     VARCHAR(36),
  removed_at     TIMESTAMP,
  removed_by     VARCHAR(36),

  CONSTRAINT pk_job PRIMARY KEY (id_job),
  CONSTRAINT fk_job_user FOREIGN KEY (id_user) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_step FOREIGN KEY (id_step) REFERENCES donation_step(id_donation_step),
  CONSTRAINT fk_job_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_job_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- CONSENT LOG
--------------------------------------------------------------------------------
CREATE TABLE consent_log (
  id_consent_log VARCHAR(36) NOT NULL,

  id_user        VARCHAR(36) NOT NULL,

  terms_version  VARCHAR(50) NOT NULL,
  accepted_at    TIMESTAMP   NOT NULL,
  ip_address     VARCHAR(45) NOT NULL,
  user_agent     TEXT        NOT NULL,

  CONSTRAINT pk_consent_log PRIMARY KEY (id_consent_log),
  CONSTRAINT uq_consent_log_user_version UNIQUE (id_user, terms_version),
  CONSTRAINT fk_consent_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- USER BABY
--------------------------------------------------------------------------------
CREATE TABLE user_baby (
  id_user_baby VARCHAR(36) NOT NULL,

  id_user      VARCHAR(36) NOT NULL,

  name         VARCHAR(120),
  birth_date   TIMESTAMP NOT NULL,

  created_at   TIMESTAMP NOT NULL,
  updated_at   TIMESTAMP,
  removed_at   TIMESTAMP,

  CONSTRAINT pk_user_baby PRIMARY KEY (id_user_baby),
  CONSTRAINT fk_user_baby_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- KB CHUNKS (base de conhecimento do módulo de IA)
--------------------------------------------------------------------------------
CREATE TABLE kb_chunks (
  id         UUID          NOT NULL DEFAULT gen_random_uuid(),

  source     VARCHAR(100)  NOT NULL,
  content    TEXT          NOT NULL,
  embedding  VECTOR(384)   NOT NULL,
  metadata   JSONB,

  created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),

  CONSTRAINT pk_kb_chunks PRIMARY KEY (id)
);

--------------------------------------------------------------------------------
-- CONVERSATIONS
--------------------------------------------------------------------------------
CREATE TABLE conversations (
  id               UUID        NOT NULL DEFAULT gen_random_uuid(),

  user_id          VARCHAR(36) NOT NULL,
  started_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_message_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  summary          TEXT,

  CONSTRAINT pk_conversations PRIMARY KEY (id),
  CONSTRAINT fk_conversations_user FOREIGN KEY (user_id) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- MESSAGES
--------------------------------------------------------------------------------
CREATE TABLE messages (
  id              UUID        NOT NULL DEFAULT gen_random_uuid(),

  conversation_id UUID        NOT NULL,
  role            VARCHAR(20) NOT NULL,
  content         TEXT        NOT NULL,
  tokens_used     INTEGER,

  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT pk_messages PRIMARY KEY (id),
  CONSTRAINT ck_messages_role CHECK (role IN ('user', 'assistant', 'system')),
  CONSTRAINT fk_messages_conversation FOREIGN KEY (conversation_id)
    REFERENCES conversations(id) ON DELETE CASCADE
);

--------------------------------------------------------------------------------
-- LLM AUDIT
--------------------------------------------------------------------------------
CREATE TABLE llm_audit (
  id              UUID        NOT NULL DEFAULT gen_random_uuid(),

  user_id         VARCHAR(36),
  conversation_id UUID,
  message_id      UUID,
  prompt_full     JSONB       NOT NULL,
  chunks_used     JSONB,
  llm_provider    VARCHAR(30) NOT NULL,
  llm_model       VARCHAR(50) NOT NULL,
  tokens_input    INTEGER,
  tokens_output   INTEGER,
  latency_ms      INTEGER,

  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  is_anonymous    BOOLEAN     NOT NULL DEFAULT false,
  session_id      VARCHAR(36),
  ip_hash         VARCHAR(64),
  action_emitted  VARCHAR(30),

  CONSTRAINT pk_llm_audit PRIMARY KEY (id)
);

--------------------------------------------------------------------------------
-- ALEMBIC VERSION (controle de migrations do serviço de IA)
--------------------------------------------------------------------------------
CREATE TABLE alembic_version (
  version_num VARCHAR(32) NOT NULL,

  CONSTRAINT pk_alembic_version PRIMARY KEY (version_num)
);

--------------------------------------------------------------------------------
-- ÍNDICES
--------------------------------------------------------------------------------
CREATE INDEX IF NOT EXISTS ix_kb_chunks_embedding
  ON kb_chunks USING hnsw (embedding vector_cosine_ops)
  WITH (m = 16, ef_construction = 64);

CREATE INDEX IF NOT EXISTS ix_conversations_user_id ON conversations (user_id);
CREATE INDEX IF NOT EXISTS ix_conversations_last_message_at ON conversations (last_message_at);
CREATE INDEX IF NOT EXISTS ix_messages_conversation_id ON messages (conversation_id);
CREATE INDEX IF NOT EXISTS ix_messages_created_at ON messages (created_at);
CREATE INDEX IF NOT EXISTS ix_llm_audit_session_id ON llm_audit (session_id);



--Massa de dados (DML):

--------------------------------------------------------------------------------
-- USER (2 adm, 3 nurse, 10 common)
--------------------------------------------------------------------------------
INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_Go8NcUY63QHtDP9bdE2zBRgFT6C', 'adm', 'Carla Menezes', '12345678901', TIMESTAMP '1985-04-12 00:00:00', '11987650001', 'carla.menezes@nutriz.org', '$2b$12$placeholderHash0000000000000000000001', TIMESTAMP '2025-01-05 09:00:00', 'system');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_TT7SOlYxGohmYipYFbxJKxDZJiN', 'adm', 'Bruno Ferreira', '12345678902', TIMESTAMP '1988-09-23 00:00:00', '11987650002', 'bruno.ferreira@nutriz.org', '$2b$12$placeholderHash0000000000000000000002', TIMESTAMP '2025-01-05 09:15:00', 'usr_Go8NcUY63QHtDP9bdE2zBRgFT6C');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', 'nurse', 'Juliana Ramos', '12345678903', TIMESTAMP '1990-01-30 00:00:00', '11987650003', 'juliana.ramos@nutriz.org', '$2b$12$placeholderHash0000000000000000000003', TIMESTAMP '2025-01-10 10:00:00', 'usr_Go8NcUY63QHtDP9bdE2zBRgFT6C');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_54TXtShO68sxNoo9iEjDVMxASJ6', 'nurse', 'Patrícia Gomes', '12345678904', TIMESTAMP '1992-06-18 00:00:00', '11987650004', 'patricia.gomes@nutriz.org', '$2b$12$placeholderHash0000000000000000000004', TIMESTAMP '2025-01-10 10:15:00', 'usr_Go8NcUY63QHtDP9bdE2zBRgFT6C');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', 'nurse', 'Fernanda Duarte', '12345678905', TIMESTAMP '1991-11-05 00:00:00', '11987650005', 'fernanda.duarte@nutriz.org', '$2b$12$placeholderHash0000000000000000000005', TIMESTAMP '2025-01-12 11:00:00', 'usr_TT7SOlYxGohmYipYFbxJKxDZJiN');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_VIZM5oForBFbyvQRZzUk1D6iNIb', 'common', 'Maria Silva', '12345678906', TIMESTAMP '1994-02-14 00:00:00', '11987650006', 'maria.silva@gmail.com', '$2b$12$placeholderHash0000000000000000000006', 1450, TIMESTAMP '2025-02-01 08:30:00', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_R13KL8bWR0rKcWWlEHPC6rlLBO0', 'common', 'Joana Souza', '12345678907', TIMESTAMP '1996-07-22 00:00:00', '11987650007', 'joana.souza@gmail.com', '$2b$12$placeholderHash0000000000000000000007', 600, TIMESTAMP '2025-02-03 09:10:00', 'usr_R13KL8bWR0rKcWWlEHPC6rlLBO0');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_u3TuU3wziWxGJgupDhrCpjgds8y', 'common', 'Camila Lima', '12345678908', TIMESTAMP '1993-03-09 00:00:00', '11987650008', 'camila.lima@gmail.com', '$2b$12$placeholderHash0000000000000000000008', TIMESTAMP '2025-02-10 14:20:00', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd', 'common', 'Beatriz Alves', '12345678909', TIMESTAMP '1990-12-01 00:00:00', '11987650009', 'beatriz.alves@gmail.com', '$2b$12$placeholderHash0000000000000000000009', 1200, TIMESTAMP '2025-02-15 16:00:00', 'usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_5GmKIjku2HChRnTLFf5GCZdDiGA', 'common', 'Larissa Costa', '12345678910', TIMESTAMP '1995-05-17 00:00:00', '11987650010', 'larissa.costa@gmail.com', '$2b$12$placeholderHash0000000000000000000010', TIMESTAMP '2025-03-01 08:00:00', 'usr_5GmKIjku2HChRnTLFf5GCZdDiGA');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_tNZPHaQ0Jt7Qg84iqh4gVJjrsMn', 'common', 'Paula Martins', '12345678911', TIMESTAMP '1989-08-25 00:00:00', '11987650011', 'paula.martins@gmail.com', '$2b$12$placeholderHash0000000000000000000011', 450, TIMESTAMP '2025-03-05 09:45:00', 'usr_tNZPHaQ0Jt7Qg84iqh4gVJjrsMn');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_5Giu7NukzNV1tLvG1GIFTKtE0bx', 'common', 'Amanda Rocha', '12345678912', TIMESTAMP '1997-10-03 00:00:00', '11987650012', 'amanda.rocha@gmail.com', '$2b$12$placeholderHash0000000000000000000012', TIMESTAMP '2025-03-10 11:30:00', 'usr_5Giu7NukzNV1tLvG1GIFTKtE0bx');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_zGD9pnLwddsFM41PREsIa2gBi4q', 'common', 'Gabriela Teixeira', '12345678913', TIMESTAMP '1992-01-19 00:00:00', '11987650013', 'gabriela.teixeira@gmail.com', '$2b$12$placeholderHash0000000000000000000013', 900, TIMESTAMP '2025-03-15 13:00:00', 'usr_zGD9pnLwddsFM41PREsIa2gBi4q');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_3IcNepOR6soVfBgWOT3gCkSt5bc', 'common', 'Rafaela Pinto', '12345678914', TIMESTAMP '1994-09-11 00:00:00', '11987650014', 'rafaela.pinto@gmail.com', '$2b$12$placeholderHash0000000000000000000014', TIMESTAMP '2025-03-20 15:15:00', 'usr_3IcNepOR6soVfBgWOT3gCkSt5bc');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_sKRWmlN2O5z0BGufzQgliEu7paq', 'common', 'Débora Nunes', '12345678915', TIMESTAMP '1991-04-28 00:00:00', '11987650015', 'debora.nunes@gmail.com', '$2b$12$placeholderHash0000000000000000000015', TIMESTAMP '2025-04-02 10:00:00', 'usr_sKRWmlN2O5z0BGufzQgliEu7paq');

--------------------------------------------------------------------------------
-- DONATION POINT (bancos/postos de leite humano reais de São Paulo)
--------------------------------------------------------------------------------
INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_fbqn1jJ8UMHBhXspthdpAOYNDeh', 'Posto de Coleta de Leite Humano Amparo Maternal', 'Posto de coleta de leite humano localizado no Amparo Maternal, Vila Clementino.', FALSE, '1150898277', 'ouvidoria@ham.spdmpais.org.br', 'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 17:00');

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_qJ38aRUhR4IWrXPvhsBkDa9U4Uq', 'Banco de Leite Humano Maternidade São Luiz Star', 'Banco de leite humano localizado na Maternidade São Luiz Star, Vila Olímpia.', FALSE, '1121211349', 'consultoriamaternidade@maternidadesaoluiz.com.br', 'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00');

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_DA4OKmTSyFzpjPSa5W3X4gXBolZ', 'Banco de Leite Humano Rede Dor São Luiz - Unidade Anália Franco', 'Banco de leite humano localizado no Hospital e Maternidade São Luiz, Unidade Anália Franco, Tatuapé.', FALSE, '1133861315', 'faleconosco.sadt@saoluiz.com.br', 'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00');

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_3PjeEYqvNSzPf2v2R2Iy9uOT4WF', 'Banco de Leite Humano do Hospital Ipiranga', 'Banco de leite humano localizado no Hospital Ipiranga, 8º andar.', TRUE, '1120677866', 'hi.cemed@gmail.com', NULL);

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_HA9k0mZMi3qdPE3xJ7gT2H2hsfW', 'Banco de Leite Humano da Santa Casa de São Paulo', 'Banco de leite humano localizado no Hospital Central da Santa Casa de São Paulo, Vila Buarque.', FALSE, '1121767390', 'ouvidoria@santacasasp.org.br', 'Seg-Sex, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00');

--------------------------------------------------------------------------------
-- ADDRESS (5 dos postos de doação + 10 residências das doadoras)
--------------------------------------------------------------------------------
INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_ren68K4TunPFz46PDjqipVJIqVL', 'dpt_fbqn1jJ8UMHBhXspthdpAOYNDeh', '04040033', 'Rua Loefgren', '101', 'São Paulo', 'SP', 'Vila Clementino', -23.600032, -46.643398, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_UxWzxczdKJmxJseyGCWJr0NRNhi', 'dpt_qJ38aRUhR4IWrXPvhsBkDa9U4Uq', '04552050', 'Rua Helena', '29', 'São Paulo', 'SP', 'Vila Olímpia', -23.590756, -46.673653, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_hP1jF7TsGTrA1EEpDJjym6MGV4i', 'dpt_DA4OKmTSyFzpjPSa5W3X4gXBolZ', '03313000', 'Rua Francisco Marengo', '1312', 'São Paulo', 'SP', 'Tatuapé', -23.548886, -46.558248, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, complement, latitude, longitude, created_at)
VALUES ('adr_GB16JC5kV3ECqWp1OrXXHFOprCe', 'dpt_3PjeEYqvNSzPf2v2R2Iy9uOT4WF', '04262000', 'Avenida Nazaré', '28', 'São Paulo', 'SP', 'Ipiranga', '8º andar', -23.584210, -46.611595, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_kxtUuX8KMf4djkWNdRfrCQBFMCA', 'dpt_HA9k0mZMi3qdPE3xJ7gT2H2hsfW', '01221020', 'Rua Dr. Cesário Mota Júnior', '112', 'São Paulo', 'SP', 'Vila Buarque', -23.542626, -46.650065, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_kplJekaACSMEscosTsS3DeRo7qY', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', '05014010', 'Rua Turiassu', '500', 'São Paulo', 'SP', 'Perdizes', -23.535000, -46.682200, TIMESTAMP '2025-02-01 08:35:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_RPPgdzUvZ3gpmmICiBlrDp37eCZ', 'usr_R13KL8bWR0rKcWWlEHPC6rlLBO0', '04101000', 'Rua Vergueiro', '1200', 'São Paulo', 'SP', 'Vila Mariana', -23.589000, -46.633900, TIMESTAMP '2025-02-03 09:15:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_41l76c1zYFl7V937s4catKMg7vs', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y', '03164000', 'Rua Coelho Barradas', '75', 'São Paulo', 'SP', 'Tatuapé', -23.540000, -46.575000, TIMESTAMP '2025-02-10 14:25:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_4fetzTUEHAXA0KeiuPeCDRHwi41', 'usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd', '04275000', 'Avenida do Cursino', '320', 'São Paulo', 'SP', 'Ipiranga', -23.592000, -46.610000, TIMESTAMP '2025-02-15 16:05:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_qD2tmy2EgpyKwKsSsb1QzraK3RX', 'usr_5GmKIjku2HChRnTLFf5GCZdDiGA', '01223000', 'Rua Barra Funda', '410', 'São Paulo', 'SP', 'Barra Funda', -23.527000, -46.662000, TIMESTAMP '2025-03-01 08:05:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_GWlG6g3Ot1OGMmjxWkI9X7H6aMu', 'usr_tNZPHaQ0Jt7Qg84iqh4gVJjrsMn', '05407000', 'Rua Cardeal Arcoverde', '1500', 'São Paulo', 'SP', 'Pinheiros', -23.567000, -46.691000, TIMESTAMP '2025-03-05 09:50:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_LdwGetDCcdx1seP32fNMGyDLJ9Y', 'usr_5Giu7NukzNV1tLvG1GIFTKtE0bx', '04547000', 'Rua Funchal', '160', 'São Paulo', 'SP', 'Vila Olímpia', -23.595000, -46.689000, TIMESTAMP '2025-03-10 11:35:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_TvnRO2qGFq562dfOB1rcavXiOqk', 'usr_zGD9pnLwddsFM41PREsIa2gBi4q', '03310000', 'Rua Serra de Botucatu', '850', 'São Paulo', 'SP', 'Tatuapé', -23.541000, -46.568000, TIMESTAMP '2025-03-15 13:05:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_YOLQZ7mBhIoPj6r0jedkYtMV0K6', 'usr_3IcNepOR6soVfBgWOT3gCkSt5bc', '04261000', 'Avenida Nazaré', '900', 'São Paulo', 'SP', 'Ipiranga', -23.586000, -46.610000, TIMESTAMP '2025-03-20 15:20:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_1Aq1fEbVId8woPeX9PcWb8pm1bN', 'usr_sKRWmlN2O5z0BGufzQgliEu7paq', '01310000', 'Avenida Paulista', '2000', 'São Paulo', 'SP', 'Bela Vista', -23.561000, -46.656000, TIMESTAMP '2025-04-02 10:05:00');

--------------------------------------------------------------------------------
-- DONATION (10: 5 concluídas, 5 em andamento)
--------------------------------------------------------------------------------
INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('don_3PvrhZxB4zVC59yvlFSFx7ZHrZf', 850.50, FALSE, 'Processo tranquilo, equipe muito atenciosa.', 5, TIMESTAMP '2025-11-10 09:00:00', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', TIMESTAMP '2025-11-20 17:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('don_XJOLlX9iBG63d1hHjtkku7Tow88', TRUE, TIMESTAMP '2026-07-05 09:00:00', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('don_VCJTBJahe84S5jIc1xLJjBictx5', 600.00, FALSE, 'Gostei muito de poder ajudar bebês prematuros.', 4, TIMESTAMP '2025-11-15 10:00:00', 'usr_R13KL8bWR0rKcWWlEHPC6rlLBO0', TIMESTAMP '2025-11-25 16:30:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('don_rJaH8mfpUAFJWpSEPTFCYbfsozS', TRUE, TIMESTAMP '2026-06-20 11:00:00', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('don_kHOUNX1yj0RpcK8ShmbCuAjASnA', 1200.00, FALSE, 'Coleta em domicílio facilitou muito.', 5, TIMESTAMP '2025-12-01 08:30:00', 'usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd', TIMESTAMP '2025-12-10 18:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('don_CciepXPxxy8KcMjRC8xxCWeKiHx', TRUE, TIMESTAMP '2026-07-10 09:30:00', 'usr_5GmKIjku2HChRnTLFf5GCZdDiGA');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('don_V5cC6ZKPmuMEGj9dCgZ51vTfGPl', 450.25, FALSE, 'Poderia ser um pouco mais rápido, mas ok.', 3, TIMESTAMP '2025-12-05 13:00:00', 'usr_tNZPHaQ0Jt7Qg84iqh4gVJjrsMn', TIMESTAMP '2025-12-18 12:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('don_uYdswxBjpHAKRYlklfN3yNRpF6L', TRUE, TIMESTAMP '2026-07-15 14:00:00', 'usr_5Giu7NukzNV1tLvG1GIFTKtE0bx');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('don_yQY0PjF9ciG9Lv3g32CgH6DaUjA', 900.00, FALSE, 'Excelente atendimento do início ao fim.', 5, TIMESTAMP '2025-12-10 09:00:00', 'usr_zGD9pnLwddsFM41PREsIa2gBi4q', TIMESTAMP '2025-12-22 17:30:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('don_7Y3c5wnRpQgwXJ43ANVj77p3kZZ', TRUE, TIMESTAMP '2026-07-20 10:00:00', 'usr_3IcNepOR6soVfBgWOT3gCkSt5bc');

--------------------------------------------------------------------------------
-- DONATION STEP (30: 4 etapas para cada doação concluída, 2 para as em andamento)
--------------------------------------------------------------------------------
-- Maria Silva (posto: Amparo Maternal / residência própria)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_krC5Gj1Bf86o0C4w7bAdzGxpyfx', 'don_3PvrhZxB4zVC59yvlFSFx7ZHrZf', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-11-11 09:00:00', TIMESTAMP '2025-11-10 09:10:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-11 10:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-11 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_iWqsuhaFVBliyIToGJ1QZwez3Vc', 'don_3PvrhZxB4zVC59yvlFSFx7ZHrZf', 'adr_kplJekaACSMEscosTsS3DeRo7qY', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-11-13 10:00:00', TIMESTAMP '2025-11-11 10:05:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-13 10:30:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-13 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_9jR7Eef1ffBgVVxZiJdL9JJvQhA', 'don_3PvrhZxB4zVC59yvlFSFx7ZHrZf', 'adr_kplJekaACSMEscosTsS3DeRo7qY', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-11-18 09:00:00', TIMESTAMP '2025-11-13 10:35:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-18 09:40:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-18 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_obu7g1TPvYjics61ES1iWTECNa5', 'don_3PvrhZxB4zVC59yvlFSFx7ZHrZf', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-11-20 15:00:00', TIMESTAMP '2025-11-18 09:45:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-20 17:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-20 17:00:00');

-- Joana Souza (posto: São Luiz Star / residência própria)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_l4AblV7vY7AZQ3VZprkYSgy3c2E', 'don_VCJTBJahe84S5jIc1xLJjBictx5', 'adr_UxWzxczdKJmxJseyGCWJr0NRNhi', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-11-16 09:00:00', TIMESTAMP '2025-11-15 10:10:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-16 10:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-16 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_Tsprvu5IfijoySjTneAAvIDAdn1', 'don_VCJTBJahe84S5jIc1xLJjBictx5', 'adr_RPPgdzUvZ3gpmmICiBlrDp37eCZ', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-11-18 10:00:00', TIMESTAMP '2025-11-16 10:05:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-18 10:30:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-18 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_4OIsP9tEpZZRztDeSdkCAEDnvMj', 'don_VCJTBJahe84S5jIc1xLJjBictx5', 'adr_RPPgdzUvZ3gpmmICiBlrDp37eCZ', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-11-22 09:00:00', TIMESTAMP '2025-11-18 10:35:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-22 09:40:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-22 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_2PP4XJU3nBC4oAv0DzAUguBuQqx', 'don_VCJTBJahe84S5jIc1xLJjBictx5', 'adr_UxWzxczdKJmxJseyGCWJr0NRNhi', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-11-24 15:00:00', TIMESTAMP '2025-11-22 09:45:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-25 16:30:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2025-11-25 16:30:00');

-- Beatriz Alves (posto: Hospital Ipiranga, com coleta domiciliar / residência própria)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_FnIiU74KKEpYEZAmggQBwBAD3Ud', 'don_kHOUNX1yj0RpcK8ShmbCuAjASnA', 'adr_GB16JC5kV3ECqWp1OrXXHFOprCe', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-12-02 09:00:00', TIMESTAMP '2025-12-01 08:40:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-02 10:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-02 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_3erXY2Av7YGr0asUt1LLQF3jCIE', 'don_kHOUNX1yj0RpcK8ShmbCuAjASnA', 'adr_4fetzTUEHAXA0KeiuPeCDRHwi41', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-12-03 10:00:00', TIMESTAMP '2025-12-02 10:05:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-03 10:30:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-03 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_rnWGhwBhsRRLFHQtcozMdant8nX', 'don_kHOUNX1yj0RpcK8ShmbCuAjASnA', 'adr_4fetzTUEHAXA0KeiuPeCDRHwi41', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-12-08 09:00:00', TIMESTAMP '2025-12-03 10:35:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-08 09:40:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-08 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_32JgdPI1af7W2pkAFEn3z5dkyay', 'don_kHOUNX1yj0RpcK8ShmbCuAjASnA', 'adr_GB16JC5kV3ECqWp1OrXXHFOprCe', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-12-10 15:00:00', TIMESTAMP '2025-12-08 09:45:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-10 18:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2025-12-10 18:00:00');

-- Paula Martins (posto: Amparo Maternal / residência própria)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_EEPmW9susPd6XfPKoIVU27c66lA', 'don_V5cC6ZKPmuMEGj9dCgZ51vTfGPl', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-12-06 09:00:00', TIMESTAMP '2025-12-05 13:10:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-06 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-06 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_6zLKQbfPBi3DldqyunDuvW4yrW8', 'don_V5cC6ZKPmuMEGj9dCgZ51vTfGPl', 'adr_GWlG6g3Ot1OGMmjxWkI9X7H6aMu', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-12-08 10:00:00', TIMESTAMP '2025-12-06 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-08 10:30:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-08 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_FfEwAvuQg2kvASFsQ8z0WJcDfuq', 'don_V5cC6ZKPmuMEGj9dCgZ51vTfGPl', 'adr_GWlG6g3Ot1OGMmjxWkI9X7H6aMu', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-12-14 09:00:00', TIMESTAMP '2025-12-08 10:35:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-14 09:40:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-14 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_joDOqDqQa5ZD5sRIkeC8wLtO9BS', 'don_V5cC6ZKPmuMEGj9dCgZ51vTfGPl', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-12-17 15:00:00', TIMESTAMP '2025-12-14 09:45:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-18 12:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-18 12:00:00');

-- Gabriela Teixeira (posto: Anália Franco / residência própria)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_zuPrp9hbVlFHy6JhqXqTCnNsS6F', 'don_yQY0PjF9ciG9Lv3g32CgH6DaUjA', 'adr_hP1jF7TsGTrA1EEpDJjym6MGV4i', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-12-11 09:00:00', TIMESTAMP '2025-12-10 09:10:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-11 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-11 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_s9Tfw7CPqVEnm0Ir7JSrigNVLpp', 'don_yQY0PjF9ciG9Lv3g32CgH6DaUjA', 'adr_TvnRO2qGFq562dfOB1rcavXiOqk', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-12-13 10:00:00', TIMESTAMP '2025-12-11 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-13 10:30:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-13 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_15FMIbOGKpTjsBaNwpKAlQQfHxe', 'don_yQY0PjF9ciG9Lv3g32CgH6DaUjA', 'adr_TvnRO2qGFq562dfOB1rcavXiOqk', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-12-18 09:00:00', TIMESTAMP '2025-12-13 10:35:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-18 09:40:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-18 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_o3RURz92ZJxfYzaqIhDxRVRqLy0', 'don_yQY0PjF9ciG9Lv3g32CgH6DaUjA', 'adr_hP1jF7TsGTrA1EEpDJjym6MGV4i', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-12-21 15:00:00', TIMESTAMP '2025-12-18 09:45:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-22 17:30:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2025-12-22 17:30:00');

-- doações em andamento: exame concluído, kit pendente
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_ptQLxEJHwBVJvwSDrtqohUmuhVI', 'don_XJOLlX9iBG63d1hHjtkku7Tow88', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-06 09:00:00', TIMESTAMP '2026-07-05 09:10:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2026-07-06 10:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', TIMESTAMP '2026-07-06 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('dst_gw2wMqZcUDIh7yfJs1ON43xKmTe', 'don_XJOLlX9iBG63d1hHjtkku7Tow88', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-06 10:05:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_jpiEQhK8nDSqXxkMM9VThX0k9tg', 'don_rJaH8mfpUAFJWpSEPTFCYbfsozS', 'adr_hP1jF7TsGTrA1EEpDJjym6MGV4i', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-06-21 09:00:00', TIMESTAMP '2026-06-20 11:10:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2026-06-21 10:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2026-06-21 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('dst_TFuEwJ77YUrshKRIy5z0w9XZjsc', 'don_rJaH8mfpUAFJWpSEPTFCYbfsozS', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-06-21 10:05:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_WJnLn03ovXjY5Mar2jiIqZlhQ3b', 'don_CciepXPxxy8KcMjRC8xxCWeKiHx', 'adr_kxtUuX8KMf4djkWNdRfrCQBFMCA', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-11 09:00:00', TIMESTAMP '2026-07-10 09:40:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2026-07-11 10:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', TIMESTAMP '2026-07-11 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('dst_TPSIAoCLrZ3aWZkSBvrjn9Wvgfy', 'don_CciepXPxxy8KcMjRC8xxCWeKiHx', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-11 10:05:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_hpoevbLJoLoaeTOdoe5c3veGprQ', 'don_uYdswxBjpHAKRYlklfN3yNRpF6L', 'adr_UxWzxczdKJmxJseyGCWJr0NRNhi', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-16 09:00:00', TIMESTAMP '2026-07-15 14:10:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2026-07-16 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2026-07-16 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('dst_tzJ1atsnBYLMPuDCCRnGEY59YVk', 'don_uYdswxBjpHAKRYlklfN3yNRpF6L', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-16 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('dst_xtl8nv8XFmoijes2YgGXI1V4HcQ', 'don_7Y3c5wnRpQgwXJ43ANVj77p3kZZ', 'adr_GB16JC5kV3ECqWp1OrXXHFOprCe', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-21 09:00:00', TIMESTAMP '2026-07-20 10:10:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2026-07-21 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', TIMESTAMP '2026-07-21 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('dst_JxhWrKoZB2JX0NNRPJbM7Q1Srbl', 'don_7Y3c5wnRpQgwXJ43ANVj77p3kZZ', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-21 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

--------------------------------------------------------------------------------
-- DONATION STEP TIMELINE (1 registro de histórico por etapa: 30)
--------------------------------------------------------------------------------
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_KmqcTBaH7ZIRU8VVQmxBe8Q6vNu', 'dst_krC5Gj1Bf86o0C4w7bAdzGxpyfx', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-11-11 10:00:00', TIMESTAMP '2025-11-11 10:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_GXN6E32VUdTiHnJuQEHyu9lD6Iv', 'dst_iWqsuhaFVBliyIToGJ1QZwez3Vc', 'adr_kplJekaACSMEscosTsS3DeRo7qY', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-11-13 10:30:00', TIMESTAMP '2025-11-13 10:30:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_om06Dwt0Y3oobQmzvr3e9XrwPGz', 'dst_9jR7Eef1ffBgVVxZiJdL9JJvQhA', 'adr_kplJekaACSMEscosTsS3DeRo7qY', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-11-18 09:40:00', TIMESTAMP '2025-11-18 09:40:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_EWIZQ0nWpRWM3YfHCHTxe6Khd1J', 'dst_obu7g1TPvYjics61ES1iWTECNa5', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-11-20 17:00:00', TIMESTAMP '2025-11-20 17:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_W7eaDNKgeInGqi7w4e4pxskC1IT', 'dst_l4AblV7vY7AZQ3VZprkYSgy3c2E', 'adr_UxWzxczdKJmxJseyGCWJr0NRNhi', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-11-16 10:00:00', TIMESTAMP '2025-11-16 10:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_3NAp935k0u7KUumWkFGDF4tFbf8', 'dst_Tsprvu5IfijoySjTneAAvIDAdn1', 'adr_RPPgdzUvZ3gpmmICiBlrDp37eCZ', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-11-18 10:30:00', TIMESTAMP '2025-11-18 10:30:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_rSWt6vwal3jKQzejVObfVHnyADv', 'dst_4OIsP9tEpZZRztDeSdkCAEDnvMj', 'adr_RPPgdzUvZ3gpmmICiBlrDp37eCZ', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-11-22 09:40:00', TIMESTAMP '2025-11-22 09:40:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_dQ5HoOodgAvTEgRXia9J7kAP744', 'dst_2PP4XJU3nBC4oAv0DzAUguBuQqx', 'adr_UxWzxczdKJmxJseyGCWJr0NRNhi', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-11-25 16:30:00', TIMESTAMP '2025-11-25 16:30:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_U5Mv17Kc03bfc8PXKqPnXKANObF', 'dst_FnIiU74KKEpYEZAmggQBwBAD3Ud', 'adr_GB16JC5kV3ECqWp1OrXXHFOprCe', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-12-02 10:00:00', TIMESTAMP '2025-12-02 10:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_BbD6e3uKBKzTOAshzb9ukZ8ND1S', 'dst_3erXY2Av7YGr0asUt1LLQF3jCIE', 'adr_4fetzTUEHAXA0KeiuPeCDRHwi41', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-12-03 10:30:00', TIMESTAMP '2025-12-03 10:30:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_ypCWr9vtLUKaqPxSpdQhDtkzRG7', 'dst_rnWGhwBhsRRLFHQtcozMdant8nX', 'adr_4fetzTUEHAXA0KeiuPeCDRHwi41', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-12-08 09:40:00', TIMESTAMP '2025-12-08 09:40:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_H5s2fqmO9JriOtNIfGPkL8LjkQN', 'dst_32JgdPI1af7W2pkAFEn3z5dkyay', 'adr_GB16JC5kV3ECqWp1OrXXHFOprCe', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-12-10 18:00:00', TIMESTAMP '2025-12-10 18:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_TjBlUH9PrNZ6IXEDB0ULru2p17f', 'dst_EEPmW9susPd6XfPKoIVU27c66lA', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-12-06 10:00:00', TIMESTAMP '2025-12-06 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_UBfBM0lIsugfuQstCMTBkSCwCcU', 'dst_6zLKQbfPBi3DldqyunDuvW4yrW8', 'adr_GWlG6g3Ot1OGMmjxWkI9X7H6aMu', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-12-08 10:30:00', TIMESTAMP '2025-12-08 10:30:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_e5fuMjeirNOLJTuyMHsDGMBgYSh', 'dst_FfEwAvuQg2kvASFsQ8z0WJcDfuq', 'adr_GWlG6g3Ot1OGMmjxWkI9X7H6aMu', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-12-14 09:40:00', TIMESTAMP '2025-12-14 09:40:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_LMcNfAQLKHu7qnQTupqziQPtDu7', 'dst_joDOqDqQa5ZD5sRIkeC8wLtO9BS', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-12-18 12:00:00', TIMESTAMP '2025-12-18 12:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_JRysqU2Q96M3jvfLQj6wt9PSQzi', 'dst_zuPrp9hbVlFHy6JhqXqTCnNsS6F', 'adr_hP1jF7TsGTrA1EEpDJjym6MGV4i', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-12-11 10:00:00', TIMESTAMP '2025-12-11 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_Vd6MVF155sXZoMZwoOmNqRWUXQR', 'dst_s9Tfw7CPqVEnm0Ir7JSrigNVLpp', 'adr_TvnRO2qGFq562dfOB1rcavXiOqk', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-12-13 10:30:00', TIMESTAMP '2025-12-13 10:30:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_R1Iv8bh4qlL9qcgMBwUYuBMGhy5', 'dst_15FMIbOGKpTjsBaNwpKAlQQfHxe', 'adr_TvnRO2qGFq562dfOB1rcavXiOqk', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-12-18 09:40:00', TIMESTAMP '2025-12-18 09:40:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_gzxYvJ8xWjmMGzGccciTvZEHDjM', 'dst_o3RURz92ZJxfYzaqIhDxRVRqLy0', 'adr_hP1jF7TsGTrA1EEpDJjym6MGV4i', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-12-22 17:30:00', TIMESTAMP '2025-12-22 17:30:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_6xfB2gpBLzHfz3tVvovXkeGOhHG', 'dst_ptQLxEJHwBVJvwSDrtqohUmuhVI', 'adr_ren68K4TunPFz46PDjqipVJIqVL', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-06 10:00:00', TIMESTAMP '2026-07-06 10:00:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('dtl_iawYYpLublqdiVAHhVeECXxGLgC', 'dst_gw2wMqZcUDIh7yfJs1ON43xKmTe', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-06 10:05:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_cQoXsf2o3gyrDO1xkxwnQrS7RPe', 'dst_jpiEQhK8nDSqXxkMM9VThX0k9tg', 'adr_hP1jF7TsGTrA1EEpDJjym6MGV4i', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-06-21 10:00:00', TIMESTAMP '2026-06-21 10:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('dtl_1iOg5OPctYCcLxUif6suVAlmiYI', 'dst_TFuEwJ77YUrshKRIy5z0w9XZjsc', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-06-21 10:05:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_E8JkqH3MB9n7IWUSmTtzQPxC5HC', 'dst_WJnLn03ovXjY5Mar2jiIqZlhQ3b', 'adr_kxtUuX8KMf4djkWNdRfrCQBFMCA', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-11 10:00:00', TIMESTAMP '2026-07-11 10:00:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('dtl_MT8ftJyPYv0iQS18VR6HfPQBGxb', 'dst_TPSIAoCLrZ3aWZkSBvrjn9Wvgfy', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-11 10:05:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_DPIHF48i2GDrmZhvkUDPqTlaVvY', 'dst_hpoevbLJoLoaeTOdoe5c3veGprQ', 'adr_UxWzxczdKJmxJseyGCWJr0NRNhi', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-16 10:00:00', TIMESTAMP '2026-07-16 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('dtl_ddLEG62Hkd9Gf2leMeR3pzh84Kp', 'dst_tzJ1atsnBYLMPuDCCRnGEY59YVk', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-16 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('dtl_Ay5XL8Sb24WKyEa8wtWy2591AIV', 'dst_xtl8nv8XFmoijes2YgGXI1V4HcQ', 'adr_GB16JC5kV3ECqWp1OrXXHFOprCe', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-21 10:00:00', TIMESTAMP '2026-07-21 10:00:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('dtl_dzrm9h820DfQnPOMbdYvpiYKne1', 'dst_JxhWrKoZB2JX0NNRPJbM7Q1Srbl', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-21 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

--------------------------------------------------------------------------------
-- JOB (10: visitas domiciliares de entrega de kit e coleta de leite)
--------------------------------------------------------------------------------
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_q7YYDsBS9UYJQTFjmsn9dLVIdVu', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', 'dst_iWqsuhaFVBliyIToGJ1QZwez3Vc', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-11-13 10:00:00', TIMESTAMP '2025-11-11 10:05:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_APg1FNA88rcSxnCC8p2xgRxI5Pw', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', 'dst_9jR7Eef1ffBgVVxZiJdL9JJvQhA', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-11-18 09:00:00', TIMESTAMP '2025-11-13 10:35:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_sWVYCoIpt9ZYE51mxR8KCDXsXyG', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', 'dst_Tsprvu5IfijoySjTneAAvIDAdn1', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-11-18 10:00:00', TIMESTAMP '2025-11-16 10:05:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_8WSlmnVErULWHMsg1msoxltaTIi', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', 'dst_4OIsP9tEpZZRztDeSdkCAEDnvMj', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-11-22 09:00:00', TIMESTAMP '2025-11-18 10:35:00', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_36wNBrOY8deQOzxGZVRk8bj2MRY', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', 'dst_3erXY2Av7YGr0asUt1LLQF3jCIE', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-12-03 10:00:00', TIMESTAMP '2025-12-02 10:05:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_IwRX3URPZSqNEm9prJto8tXsTnS', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', 'dst_rnWGhwBhsRRLFHQtcozMdant8nX', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-12-08 09:00:00', TIMESTAMP '2025-12-03 10:35:00', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_vRhALtY5U3SObMEq9PYXLKoUdLE', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', 'dst_6zLKQbfPBi3DldqyunDuvW4yrW8', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-12-08 10:00:00', TIMESTAMP '2025-12-06 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_Fbh7x41Ztpdp4K8ffUF0eWIXiiQ', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', 'dst_FfEwAvuQg2kvASFsQ8z0WJcDfuq', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-12-14 09:00:00', TIMESTAMP '2025-12-08 10:35:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_w3Q8WB36Ud9sMtwgKGnjQEo2gw2', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', 'dst_s9Tfw7CPqVEnm0Ir7JSrigNVLpp', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-12-13 10:00:00', TIMESTAMP '2025-12-11 10:05:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_9SHDdJp62hDiZDQHJMu8W5CN0U5', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', 'dst_15FMIbOGKpTjsBaNwpKAlQQfHxe', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-12-18 09:00:00', TIMESTAMP '2025-12-13 10:35:00', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4');

--------------------------------------------------------------------------------
-- CONSENT LOG (15: um por usuário)
--------------------------------------------------------------------------------
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_i7sC4SFhbOMZpTktJaJAfo16hD8', 'usr_Go8NcUY63QHtDP9bdE2zBRgFT6C', 'v1.0', TIMESTAMP '2025-01-05 09:00:00', '200.150.10.11', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_Q2hU5tGtQAuzSsJimAQ8yRV5lNK', 'usr_TT7SOlYxGohmYipYFbxJKxDZJiN', 'v1.0', TIMESTAMP '2025-01-05 09:15:00', '200.150.10.12', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_O8xgRoEbN48Ju6NoPeO0D6StPAh', 'usr_Lb7tKR69yz8TmeLS1OpgSXt2RMZ', 'v1.0', TIMESTAMP '2025-01-10 10:00:00', '200.150.10.13', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_rcdJsS8iO3WFg3aKsEECvl9dq83', 'usr_54TXtShO68sxNoo9iEjDVMxASJ6', 'v1.0', TIMESTAMP '2025-01-10 10:15:00', '200.150.10.14', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_Eh0ezFeKORdjjZK8tfphJWAMMYN', 'usr_r4CpWDKNQyvbF2ulFnwZqvr4MS4', 'v1.0', TIMESTAMP '2025-01-12 11:00:00', '200.150.10.15', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_B5LzxoiGFfWd3hjOkYRBMeyyMDH', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', 'v1.1', TIMESTAMP '2025-02-01 08:30:00', '187.10.20.1', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_cpTCCHHNkxx6syAXvRMdYOPvevg', 'usr_R13KL8bWR0rKcWWlEHPC6rlLBO0', 'v1.1', TIMESTAMP '2025-02-03 09:10:00', '187.10.20.2', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_9HIGYGJby3EcOyxqVbwYewpUQOg', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y', 'v1.1', TIMESTAMP '2025-02-10 14:20:00', '187.10.20.3', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_mhi2eCl5TCfZR92uQwTeJIs5t2k', 'usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd', 'v1.1', TIMESTAMP '2025-02-15 16:00:00', '187.10.20.4', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_wvJWIyD7u3mSpKyo2XAcuVET6Zy', 'usr_5GmKIjku2HChRnTLFf5GCZdDiGA', 'v1.1', TIMESTAMP '2025-03-01 08:00:00', '187.10.20.5', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_oXHyC6Ct3LBtKNdN9Vg8WnOnqQf', 'usr_tNZPHaQ0Jt7Qg84iqh4gVJjrsMn', 'v1.1', TIMESTAMP '2025-03-05 09:45:00', '187.10.20.6', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_QfsGQONvf08WpRtoZmjbcpEN2Xe', 'usr_5Giu7NukzNV1tLvG1GIFTKtE0bx', 'v1.1', TIMESTAMP '2025-03-10 11:30:00', '187.10.20.7', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_ic8ctFhgp4IiyDxQ8VS8IALVUj4', 'usr_zGD9pnLwddsFM41PREsIa2gBi4q', 'v1.1', TIMESTAMP '2025-03-15 13:00:00', '187.10.20.8', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_FDWKWLCR74KPONu3OujCeECOtYr', 'usr_3IcNepOR6soVfBgWOT3gCkSt5bc', 'v1.1', TIMESTAMP '2025-03-20 15:15:00', '187.10.20.9', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('clg_hYKYcwIBQxeGPva2A0FgB9xO51D', 'usr_sKRWmlN2O5z0BGufzQgliEu7paq', 'v1.1', TIMESTAMP '2025-04-02 10:00:00', '187.10.20.10', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');

--------------------------------------------------------------------------------
-- USER BABY (8: filhos das doadoras, motivo da amamentação/doação)
--------------------------------------------------------------------------------
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_XLVWvicwIv0Pl1XRDSEOlZieTX8', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', 'Sofia', TIMESTAMP '2025-01-10 00:00:00', TIMESTAMP '2025-02-01 08:32:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_sChDStSz8rGIFCfMc4BVuMqbfo9', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', 'Miguel', TIMESTAMP '2025-01-10 00:00:00', TIMESTAMP '2025-02-01 08:33:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_m5XwwU90P0jpgjqmlMjWWPel8XO', 'usr_R13KL8bWR0rKcWWlEHPC6rlLBO0', 'Helena', TIMESTAMP '2025-03-22 00:00:00', TIMESTAMP '2025-02-03 09:12:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_4xHG6r1kq608E9ZsV3vZhD9ejW9', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y', 'Davi', TIMESTAMP '2025-05-02 00:00:00', TIMESTAMP '2025-02-10 14:22:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_uhXz3G0aQ3IDAdmHxNWFOCWdnrJ', 'usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd', 'Laura', TIMESTAMP '2024-11-15 00:00:00', TIMESTAMP '2025-02-15 16:02:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_v4XNiMyjkl1SXNZ5kUCcAxRUpCN', 'usr_5GmKIjku2HChRnTLFf5GCZdDiGA', 'Pedro', TIMESTAMP '2025-02-08 00:00:00', TIMESTAMP '2025-03-01 08:02:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_OhbVrpoiVgRV5IfLBcbfnoGMbJm', 'usr_tNZPHaQ0Jt7Qg84iqh4gVJjrsMn', 'Isabela', TIMESTAMP '2025-06-19 00:00:00', TIMESTAMP '2025-03-05 09:47:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('usb_MOkIUpkDyr7OSJoRu1XXdo0cZuz', 'usr_zGD9pnLwddsFM41PREsIa2gBi4q', 'Enzo', TIMESTAMP '2025-04-30 00:00:00', TIMESTAMP '2025-03-15 13:02:00');

--------------------------------------------------------------------------------
-- KB CHUNKS (6: base de conhecimento do assistente de IA)
-- Os embeddings de 384 dimensões são gerados em PL/pgSQL (impraticável
-- digitar 384 números manualmente); o conteúdo textual é real e coerente
-- com o domínio.
--------------------------------------------------------------------------------
DO $$
DECLARE
  v_ids      UUID[] := ARRAY[
    '22f86792-1b03-4b72-badd-37ea17dd1590',
    '7d84ddb1-a9f2-4ca5-92bb-6344cdb6c7ef',
    '5d968845-5cae-4ced-a6c6-eb33cd84b093',
    'd58b8009-7559-4a5e-a8ac-321e1dce4d08',
    'ccfd8cb5-a6cd-4526-9c8d-fb5b15421de9',
    'e106098f-f104-4037-a30a-caf8d6c628f8'
  ];
  v_sources  TEXT[] := ARRAY[
    'faq_doacao_leite.pdf',
    'manual_banco_leite.pdf',
    'manual_banco_leite.pdf',
    'faq_doacao_leite.pdf',
    'guia_amamentacao.pdf',
    'guia_amamentacao.pdf'
  ];
  v_contents TEXT[] := ARRAY[
    'Podem doar leite humano mulheres saudáveis, não fumantes, que não usam determinados medicamentos e que produzem leite excedente após amamentar o próprio bebê.',
    'Antes da ordenha, a doadora deve lavar as mãos e os seios com água e sabão neutro, e utilizar frascos esterilizados fornecidos pelo banco de leite.',
    'O leite ordenhado deve ser armazenado em freezer, identificado com data e hora da coleta, podendo ser mantido por até 15 dias antes da coleta pelo banco de leite.',
    'A doadora passa por uma triagem com exame de sangue para descartar doenças infectocontagiosas antes de iniciar as doações regulares.',
    'O leite humano doado é destinado prioritariamente a bebês prematuros e internados em UTI neonatal, auxiliando no fortalecimento do sistema imunológico.',
    'O kit de ordenha, com frascos e etiquetas, é entregue gratuitamente na casa da doadora ou em um dos postos de coleta cadastrados no Nutriz.'
  ];
  v_metadata TEXT[] := ARRAY[
    '{"topic":"elegibilidade","lang":"pt-BR"}',
    '{"topic":"higiene","lang":"pt-BR"}',
    '{"topic":"armazenamento","lang":"pt-BR"}',
    '{"topic":"exames","lang":"pt-BR"}',
    '{"topic":"beneficios","lang":"pt-BR"}',
    '{"topic":"kit_coleta","lang":"pt-BR"}'
  ];
  v_vec TEXT;
BEGIN
  FOR i IN 1 .. array_length(v_ids, 1) LOOP
    v_vec := '[';
    FOR j IN 1 .. 384 LOOP
      v_vec := v_vec || round((random() * 2 - 1)::numeric, 4)::text;
      IF j < 384 THEN
        v_vec := v_vec || ',';
      END IF;
    END LOOP;
    v_vec := v_vec || ']';

    INSERT INTO kb_chunks (id, source, content, embedding, metadata)
    VALUES (
      v_ids[i],
      v_sources[i],
      v_contents[i],
      v_vec::vector,
      v_metadata[i]::jsonb
    );
  END LOOP;
END $$;

--------------------------------------------------------------------------------
-- CONVERSATIONS (8: doadoras conversando com o assistente de IA)
--------------------------------------------------------------------------------
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('bd23c48c-395d-4f4e-881e-1573083073ab', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', TIMESTAMP '2026-07-05 08:00:00', TIMESTAMP '2026-07-05 08:06:00', 'Dúvidas sobre nova doação e reagendamento de etapas.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('581c3349-7563-4f84-853c-6656df9f3f49', 'usr_R13KL8bWR0rKcWWlEHPC6rlLBO0', TIMESTAMP '2026-07-06 09:00:00', TIMESTAMP '2026-07-06 09:03:00', 'Dúvida sobre prazo de armazenamento do leite no freezer.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('40f25e2f-aa6e-4142-aa71-9d9d39caa27a', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y', TIMESTAMP '2026-06-19 10:00:00', TIMESTAMP '2026-06-19 10:08:00', 'Primeira doadora perguntando como iniciar o processo de doação.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('3d636b24-8382-4610-bc4d-9cb285992752', 'usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd', TIMESTAMP '2025-12-11 11:00:00', TIMESTAMP '2025-12-11 11:02:00', 'Consulta sobre resultado da análise de leite.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('975fdb35-d7b4-432e-b670-6b8241242a1e', 'usr_5GmKIjku2HChRnTLFf5GCZdDiGA', TIMESTAMP '2026-07-09 15:00:00', TIMESTAMP '2026-07-09 15:02:00', 'Solicitação de alteração de endereço para coleta.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('1c2d3f81-81f0-4a8a-ad28-02aee71ba779', 'usr_5Giu7NukzNV1tLvG1GIFTKtE0bx', TIMESTAMP '2026-07-16 16:00:00', TIMESTAMP '2026-07-16 16:02:00', 'Dúvida sobre prazo de entrega do kit após exame.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('c606df53-03d2-4525-a3b8-d3238a041b7f', 'usr_3IcNepOR6soVfBgWOT3gCkSt5bc', TIMESTAMP '2026-07-20 17:00:00', TIMESTAMP '2026-07-20 17:02:00', 'Pergunta sobre coleta domiciliar no Hospital Ipiranga.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('a96c46f9-118f-4836-9cef-3ee3b8e8b54b', 'usr_sKRWmlN2O5z0BGufzQgliEu7paq', TIMESTAMP '2026-08-01 12:00:00', TIMESTAMP '2026-08-01 12:03:00', 'Usuária ainda não doadora perguntando sobre requisitos de elegibilidade.');

--------------------------------------------------------------------------------
-- MESSAGES (20: alternando usuário/assistente)
--------------------------------------------------------------------------------
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('1df353e3-669e-490d-8279-1fdd234bc1d0', 'bd23c48c-395d-4f4e-881e-1573083073ab', 'user', 'Já fiz uma doação em novembro, posso doar de novo?', 14, TIMESTAMP '2026-07-05 08:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('01397565-4f36-4e08-9a25-4ebc02ad159e', 'bd23c48c-395d-4f4e-881e-1573083073ab', 'assistant', 'Sim! Você pode doar novamente após passar por uma nova avaliação de saúde. Vou agendar o exame de sangue para você.', 28, TIMESTAMP '2026-07-05 08:02:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('ce4d76df-2dc7-4bb8-9b71-a626abb983f8', 'bd23c48c-395d-4f4e-881e-1573083073ab', 'user', 'Perfeito, quando será a entrega do kit?', 10, TIMESTAMP '2026-07-05 08:05:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('2e2454e8-c091-44a0-bba1-92fdee3117bc', 'bd23c48c-395d-4f4e-881e-1573083073ab', 'assistant', 'A entrega do kit de ordenha está aguardando agendamento; uma de nossas enfermeiras entrará em contato para confirmar o horário.', 26, TIMESTAMP '2026-07-05 08:06:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('034e9aa0-efe0-4cdf-b539-eb1a56bf139c', '581c3349-7563-4f84-853c-6656df9f3f49', 'user', 'Quanto tempo posso guardar o leite ordenhado no freezer antes da coleta?', 16, TIMESTAMP '2026-07-06 09:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('4863d92d-6221-45eb-ad6b-468631f755d4', '581c3349-7563-4f84-853c-6656df9f3f49', 'assistant', 'O leite pode ser armazenado no freezer por até 15 dias antes da coleta pelo banco de leite, sempre identificado com data e hora.', 24, TIMESTAMP '2026-07-06 09:03:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('06f2cd1b-2d94-4c4d-bb33-7bf7ec0a9506', '40f25e2f-aa6e-4142-aa71-9d9d39caa27a', 'user', 'Estou amamentando meu filho Davi e tenho leite excedente, como funciona a doação?', 19, TIMESTAMP '2026-06-19 10:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('e2bc3742-8b4e-4d75-9449-02f2e7c078ee', '40f25e2f-aa6e-4142-aa71-9d9d39caa27a', 'assistant', 'Que ótimo, Camila! Após um exame de sangue de triagem, você recebe um kit de ordenha e pode começar a doar. Deseja que eu inicie o processo?', 30, TIMESTAMP '2026-06-19 10:04:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('ea14acd6-85b0-411b-aff0-2c2828a41d5f', '40f25e2f-aa6e-4142-aa71-9d9d39caa27a', 'user', 'Sim, pode iniciar.', 6, TIMESTAMP '2026-06-19 10:06:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('0e4ea09f-22c8-48f2-b9ec-874728f38bb3', '40f25e2f-aa6e-4142-aa71-9d9d39caa27a', 'assistant', 'Prontinho! Já registrei sua doação e agendei o exame de sangue inicial.', 15, TIMESTAMP '2026-06-19 10:08:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('b2a9348c-845b-4dfd-bf0a-0fcb69de0289', '3d636b24-8382-4610-bc4d-9cb285992752', 'user', 'Recebi a notificação de que minha doação foi analisada, qual foi o resultado?', 17, TIMESTAMP '2025-12-11 11:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('fbbda0d1-7846-4bde-80ee-79fbda946e53', '3d636b24-8382-4610-bc4d-9cb285992752', 'assistant', 'Sua doação de 1200ml foi aprovada na análise de leite e já está disponível para os bebês atendidos pelo Hospital Ipiranga.', 27, TIMESTAMP '2025-12-11 11:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('98825c30-340a-4ae9-b77d-2de153d4f664', '975fdb35-d7b4-432e-b670-6b8241242a1e', 'user', 'Preciso trocar o endereço para a coleta do kit, como faço?', 14, TIMESTAMP '2026-07-09 15:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('82caddad-a423-4b31-922c-329e5ee8bf1e', '975fdb35-d7b4-432e-b670-6b8241242a1e', 'assistant', 'Claro, você pode atualizar seu endereço no aplicativo e a próxima etapa será ajustada automaticamente.', 20, TIMESTAMP '2026-07-09 15:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('8f385499-3652-4709-bdfb-05355aada896', '1c2d3f81-81f0-4a8a-ad28-02aee71ba779', 'user', 'Fiz o exame de sangue, quando recebo o kit de ordenha?', 13, TIMESTAMP '2026-07-16 16:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('848e9e6a-1042-4e80-8e52-96dc712588e4', '1c2d3f81-81f0-4a8a-ad28-02aee71ba779', 'assistant', 'Seu exame foi processado e o kit de ordenha será entregue em até 3 dias úteis no endereço cadastrado.', 21, TIMESTAMP '2026-07-16 16:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('3eae21b7-67b2-47a1-9e3c-b9519ba37dd5', 'c606df53-03d2-4525-a3b8-d3238a041b7f', 'user', 'O Hospital Ipiranga faz coleta em domicílio?', 10, TIMESTAMP '2026-07-20 17:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('707cfc3b-0907-4606-9e23-b4dcb710d749', 'c606df53-03d2-4525-a3b8-d3238a041b7f', 'assistant', 'Sim! O Banco de Leite Humano do Hospital Ipiranga oferece coleta domiciliar, ideal para o seu caso.', 19, TIMESTAMP '2026-07-20 17:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('695832b5-a0fa-41fb-af32-7017cba68232', 'a96c46f9-118f-4836-9cef-3ee3b8e8b54b', 'user', 'Ainda não doei, quais são os requisitos para me tornar doadora?', 15, TIMESTAMP '2026-08-01 12:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('9bf7cfdd-4d27-41b0-963f-455d15899ab7', 'a96c46f9-118f-4836-9cef-3ee3b8e8b54b', 'assistant', 'Você precisa ser uma mulher saudável, não fumante, sem uso de determinados medicamentos, e produzir leite excedente após amamentar seu bebê. Deseja iniciar a triagem?', 32, TIMESTAMP '2026-08-01 12:03:00');

--------------------------------------------------------------------------------
-- LLM AUDIT (10: uma auditoria por resposta do assistente)
--------------------------------------------------------------------------------
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('18dfb96c-0603-4df0-9776-48bc4c69cb0c', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', 'bd23c48c-395d-4f4e-881e-1573083073ab', '01397565-4f36-4e08-9a25-4ebc02ad159e', '{"system":"Assistente Nutriz","user":"Já fiz uma doação em novembro, posso doar de novo?"}', '["kb_faq_elegibilidade"]', 'anthropic', 'claude-sonnet-5', 120, 28, 850, FALSE, 'sess_maria_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('189789ce-dad8-4795-bfb5-071f88d8d307', 'usr_VIZM5oForBFbyvQRZzUk1D6iNIb', 'bd23c48c-395d-4f4e-881e-1573083073ab', '2e2454e8-c091-44a0-bba1-92fdee3117bc', '{"system":"Assistente Nutriz","user":"Perfeito, quando será a entrega do kit?"}', '["kb_faq_kit_coleta"]', 'anthropic', 'claude-sonnet-5', 95, 26, 780, FALSE, 'sess_maria_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('a3133f41-8740-4552-a348-7ae44ab832a9', 'usr_R13KL8bWR0rKcWWlEHPC6rlLBO0', '581c3349-7563-4f84-853c-6656df9f3f49', '4863d92d-6221-45eb-ad6b-468631f755d4', '{"system":"Assistente Nutriz","user":"Quanto tempo posso guardar o leite ordenhado no freezer antes da coleta?"}', '["kb_faq_armazenamento"]', 'anthropic', 'claude-sonnet-5', 110, 24, 720, FALSE, 'sess_joana_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('1bccd1b7-adb0-4c74-ae53-75307c50f26a', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y', '40f25e2f-aa6e-4142-aa71-9d9d39caa27a', 'e2bc3742-8b4e-4d75-9449-02f2e7c078ee', '{"system":"Assistente Nutriz","user":"Estou amamentando meu filho Davi e tenho leite excedente, como funciona a doação?"}', '["kb_faq_elegibilidade","kb_faq_kit_coleta"]', 'anthropic', 'claude-sonnet-5', 140, 30, 910, FALSE, 'sess_camila_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('0bba9379-f6a4-4284-a35e-13a0c68ae3a5', 'usr_u3TuU3wziWxGJgupDhrCpjgds8y', '40f25e2f-aa6e-4142-aa71-9d9d39caa27a', '0e4ea09f-22c8-48f2-b9ec-874728f38bb3', '{"system":"Assistente Nutriz","user":"Sim, pode iniciar."}', NULL, 'anthropic', 'claude-sonnet-5', 60, 15, 500, FALSE, 'sess_camila_01', 'start_donation');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('7d7f6d6d-e0f5-40bf-b3bd-1b50d49a354d', 'usr_DcsmcYm4cu7tGz0IEqcWPmsw3Xd', '3d636b24-8382-4610-bc4d-9cb285992752', 'fbbda0d1-7846-4bde-80ee-79fbda946e53', '{"system":"Assistente Nutriz","user":"Recebi a notificação de que minha doação foi analisada, qual foi o resultado?"}', NULL, 'anthropic', 'claude-sonnet-5', 105, 27, 760, FALSE, 'sess_beatriz_01', 'answer_status');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('b4e4a0f2-cf73-4342-8a15-80b974dbd096', 'usr_5GmKIjku2HChRnTLFf5GCZdDiGA', '975fdb35-d7b4-432e-b670-6b8241242a1e', '82caddad-a423-4b31-922c-329e5ee8bf1e', '{"system":"Assistente Nutriz","user":"Preciso trocar o endereço para a coleta do kit, como faço?"}', NULL, 'anthropic', 'claude-sonnet-5', 90, 20, 640, FALSE, 'sess_larissa_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('8f836fb4-6158-4850-a75e-077d77d4045a', 'usr_5Giu7NukzNV1tLvG1GIFTKtE0bx', '1c2d3f81-81f0-4a8a-ad28-02aee71ba779', '848e9e6a-1042-4e80-8e52-96dc712588e4', '{"system":"Assistente Nutriz","user":"Fiz o exame de sangue, quando recebo o kit de ordenha?"}', '["kb_faq_kit_coleta"]', 'anthropic', 'claude-sonnet-5', 100, 21, 700, FALSE, 'sess_amanda_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('befe59ea-3a73-4c1c-945c-f8d82a69439b', 'usr_3IcNepOR6soVfBgWOT3gCkSt5bc', 'c606df53-03d2-4525-a3b8-d3238a041b7f', '707cfc3b-0907-4606-9e23-b4dcb710d749', '{"system":"Assistente Nutriz","user":"O Hospital Ipiranga faz coleta em domicílio?"}', NULL, 'anthropic', 'claude-sonnet-5', 85, 19, 610, FALSE, 'sess_rafaela_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('29d4e491-d191-48b5-bbf1-f69bcc34bfb6', 'usr_sKRWmlN2O5z0BGufzQgliEu7paq', 'a96c46f9-118f-4836-9cef-3ee3b8e8b54b', '9bf7cfdd-4d27-41b0-963f-455d15899ab7', '{"system":"Assistente Nutriz","user":"Ainda não doei, quais são os requisitos para me tornar doadora?"}', '["kb_faq_elegibilidade"]', 'anthropic', 'claude-sonnet-5', 130, 32, 880, FALSE, 'sess_debora_01', 'answer_faq');

--------------------------------------------------------------------------------
-- ALEMBIC VERSION (controle de migrations do serviço de IA)
--------------------------------------------------------------------------------
INSERT INTO alembic_version (version_num) VALUES ('d4e91a7c22b0');



-- Consultas SQL:

-- Consultas e Inteligência de Negócio

-- 1) Doadoras cadastradas (usuárias comuns), ordenadas por nome.
SELECT id_user, name, email, milk_donated
FROM "user"
WHERE type = 'common'
ORDER BY name;

-- 2) Doações já concluídas e bem avaliadas (nota >= 4), da mais recente para a mais antiga.
SELECT id_donation, quantity_donated, score_feedback, created_at
FROM donation
WHERE is_active = FALSE
  AND score_feedback >= 4
ORDER BY created_at DESC;

-- 3) Bairros distintos onde há pontos de doação cadastrados.
SELECT DISTINCT a.neighborhood
FROM address a
JOIN donation_point dp ON dp.id_donation_point = a.id_donation_point
ORDER BY a.neighborhood;

-- 4) Etapas de doação ainda pendentes, ordenadas pela data prevista.
SELECT id_donation_step, id_donation, name, status, set_date
FROM donation_step
WHERE status = 'pending'
ORDER BY set_date;

-- 5) Tipos de perfil distintos cadastrados no sistema.
SELECT DISTINCT type
FROM "user"
ORDER BY type;


-- Consultas Realacionais

-- 1) INNER JOIN: doações com o nome da doadora responsável.
SELECT d.id_donation, u.name AS doadora, d.quantity_donated, d.is_active, d.created_at
FROM donation d
INNER JOIN "user" u ON u.id_user = d.created_by
ORDER BY d.created_at DESC;

-- 2) INNER JOIN (múltiplas tabelas): etapas de doação com a doadora e a enfermeira responsável.
SELECT ds.id_donation_step, ds.name AS etapa, ds.status,
       ud.name AS doadora, un.name AS enfermeira_responsavel
FROM donation_step ds
INNER JOIN donation dn ON dn.id_donation = ds.id_donation
INNER JOIN "user" ud ON ud.id_user = dn.created_by
INNER JOIN "user" un ON un.id_user = ds.created_by
ORDER BY dn.id_donation, ds.created_at;

-- 3) LEFT JOIN: todos os usuários e seus bebês cadastrados (inclui quem não tem bebê registrado).
SELECT u.name AS usuario, b.name AS bebe, b.birth_date
FROM "user" u
LEFT JOIN user_baby b ON b.id_user = u.id_user
ORDER BY u.name, b.name;

-- 4) LEFT JOIN: todos os postos de doação e o respectivo endereço (inclui postos sem endereço cadastrado).
SELECT dp.name AS posto, a.street, a."number", a.neighborhood
FROM donation_point dp
LEFT JOIN address a ON a.id_donation_point = dp.id_donation_point
ORDER BY dp.name;

-- 5) RIGHT JOIN: todos os usuários e seu endereço, mesmo quem não cadastrou um (ex.: equipe administrativa).
SELECT u.name AS usuario, a.street, a."number", a.neighborhood
FROM address a
RIGHT JOIN "user" u ON u.id_user = a.id_user
ORDER BY u.name;


-- Consultas analíticas 

-- 1) Ranking de doadoras: volume total doado e número de doações concluídas.
--    Indicador: identifica as maiores doadoras para ações de reconhecimento/fidelização.
SELECT u.name AS doadora,
       COUNT(d.id_donation) AS qtd_doacoes,
       SUM(d.quantity_donated) AS total_doado_ml
FROM donation d
INNER JOIN "user" u ON u.id_user = d.created_by
WHERE d.quantity_donated IS NOT NULL
GROUP BY u.name
ORDER BY total_doado_ml DESC;

-- 2) Doadoras recorrentes (mais de uma doação registrada, concluída ou em andamento).
--    Indicador: mede fidelização/retenção de doadoras.
SELECT u.name AS doadora, COUNT(*) AS qtd_doacoes
FROM donation d
INNER JOIN "user" u ON u.id_user = d.created_by
GROUP BY u.name
HAVING COUNT(*) > 1
ORDER BY qtd_doacoes DESC;

-- 3) Avaliação média por posto de coleta, somente postos com pelo menos 2 doações avaliadas.
--    Indicador: qualidade percebida do atendimento por unidade, com amostra mínima confiável.
SELECT dp.name AS posto,
       COUNT(DISTINCT d.id_donation) AS qtd_doacoes_avaliadas,
       ROUND(AVG(d.score_feedback), 2) AS media_avaliacao
FROM donation d
INNER JOIN donation_step ds ON ds.id_donation = d.id_donation AND ds.name = 'Análise de leite'
INNER JOIN address a ON a.id_address = ds.id_address
INNER JOIN donation_point dp ON dp.id_donation_point = a.id_donation_point
WHERE d.score_feedback IS NOT NULL
GROUP BY dp.name
HAVING COUNT(DISTINCT d.id_donation) >= 2
ORDER BY media_avaliacao DESC;

-- 4) Volume de etapas do pipeline de doação por nome e status.
--    Indicador: mostra em quais etapas há mais gargalos (ex.: muitas etapas "pending").
SELECT name AS etapa, status, COUNT(*) AS qtd
FROM donation_step
GROUP BY name, status
ORDER BY etapa, status;

-- 5) Produtividade da equipe: enfermeiras com 2 ou mais visitas (jobs) concluídas.
--    Indicador: apoia decisões de alocação de equipe para visitas domiciliares.
SELECT u.name AS enfermeira, COUNT(*) AS jobs_concluidos
FROM job j
INNER JOIN "user" u ON u.id_user = j.id_user
WHERE j.status = 'done'
GROUP BY u.name
HAVING COUNT(*) >= 2
ORDER BY jobs_concluidos DESC;
