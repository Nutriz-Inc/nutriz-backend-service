--------------------------------------------------------------------------------
-- nutriz-postgres.sql
-- DDL consolidado (PostgreSQL) do banco de dados Nutriz, gerado a partir de
-- todas as migrations em shared/database/migrations/ (0_init .. 6_remove_user_uniques).
--
-- Equivalente PostgreSQL de nutriz.sql (versão Oracle). As duas versões
-- representam o mesmo modelo lógico; as diferenças abaixo existem apenas
-- porque cada banco expõe recursos nativos distintos:
--   * ENUM (CREATE TYPE ... AS ENUM) nativo, em vez de VARCHAR2 + CHECK.
--   * BOOLEAN nativo (sempre existiu no Postgres, sem exigir versão especial).
--   * UUID nativo com gen_random_uuid() (extensão pgcrypto), em vez de uma
--     função wrapper sobre SYS_GUID().
--   * VECTOR(384) via extensão pgvector, com índice HNSW nativo.
--   * A coluna address.number não precisa de aspas: NUMBER não é palavra
--     reservada no Postgres (o tipo numérico chama-se NUMERIC).
--   * generate_ksuid() usa a função original em PL/pgSQL (aritmética de
--     128 bits), pois o NUMERIC do Postgres suporta a precisão necessária
--     sem o estouro que obrigou a simplificar a versão Oracle.
--   * "user" permanece entre aspas duplas: USER é palavra reservada também
--     no Postgres (idêntico à migration original).
--
-- Assim como na versão Oracle, foram adicionadas as constraints que a
-- migration 6 removeu de "user" (ver comentário na tabela) mais um punhado
-- de UNIQUE/CHECK de regra de negócio ausentes nas migrations originais,
-- para que o DDL final cubra PRIMARY KEY, FOREIGN KEY, UNIQUE, CHECK e
-- NOT NULL em pelo menos uma tabela cada.
--------------------------------------------------------------------------------

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
