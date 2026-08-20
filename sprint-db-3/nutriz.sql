--------------------------------------------------------------------------------
-- nutriz.sql
-- DDL consolidado (Oracle SQL) do banco de dados Nutriz, gerado a partir de
-- todas as migrations em shared/database/migrations/ (0_init .. 6_remove_user_uniques).
--
-- Notas de conversão a partir do schema original (PostgreSQL):
--   * VARCHAR/TEXT           -> VARCHAR2/CLOB
--   * BOOLEAN                -> BOOLEAN (tipo nativo SQL, disponível a partir do
--     Oracle Database 23c/23ai; mesmo requisito de versão já usado pelo VECTOR)
--   * ENUM (CREATE TYPE ...) -> VARCHAR2 + CHECK (col IN (...))
--   * UUID / gen_random_uuid -> VARCHAR2(36), preenchido via generate_uuid()
--   * JSONB                  -> CLOB + CHECK (col IS JSON)
--   * VECTOR(384)            -> VECTOR(384, FLOAT32) (requer Oracle Database 23ai)
--   * "user" permanece entre aspas duplas: USER é palavra reservada no Oracle
--     (pseudocoluna da sessão corrente); toda referência à tabela deve manter
--     a mesma grafia entre aspas.
--   * A coluna address."number" é mantida entre aspas: NUMBER é palavra
--     reservada do Oracle (nome de tipo de dado).
--   * Nomes de constraints de FOREIGN KEY duplicados entre tabelas nas
--     migrations originais (ex.: fk_donation_created_by reaproveitado em
--     donation, donation_step e donation_step_timeline) foram tornados únicos,
--     pois no Oracle o nome de uma constraint deve ser único em todo o schema.
--   * Migrations de dados (INSERT/UPDATE/DELETE) foram ignoradas: apenas DDL.
--   * A migration 6 remove as UNIQUE constraints de "user" (internal_identifier,
--     cpf, phone_number, email); portanto o estado final de "user" não as
--     possui. As colunas continuam NOT NULL onde a migration original definia.
--------------------------------------------------------------------------------

--------------------------------------------------------------------------------
-- FUNÇÕES AUXILIARES
--------------------------------------------------------------------------------

-- Gera um identificador UUID-like (formato 8-4-4-4-12), substituindo o uso de
-- gen_random_uuid()/pgcrypto do PostgreSQL para as tabelas do módulo de IA.
CREATE OR REPLACE FUNCTION generate_uuid RETURN VARCHAR2 IS
  v_hex VARCHAR2(32) := RAWTOHEX(SYS_GUID());
BEGIN
  RETURN LOWER(
    SUBSTR(v_hex, 1, 8)  || '-' ||
    SUBSTR(v_hex, 9, 4)  || '-' ||
    SUBSTR(v_hex, 13, 4) || '-' ||
    SUBSTR(v_hex, 17, 4) || '-' ||
    SUBSTR(v_hex, 21, 12)
  );
END generate_uuid;
/

-- Gera um identificador base62 de 27 caracteres, no mesmo formato usado pelos
-- ids prefixados de donation_point e address (ver migration 2_update_default_id).
-- Simplificação da função generate_ksuid() original em plpgsql: aqui o valor
-- é puramente aleatório (não time-sortable), pois a aritmética de 128 bits do
-- KSUID original excede a precisão máxima (38 dígitos) do tipo NUMBER do Oracle.
CREATE OR REPLACE FUNCTION generate_ksuid RETURN VARCHAR2 IS
  v_alphabet VARCHAR2(62) := '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
  v_result   VARCHAR2(27) := '';
BEGIN
  FOR i IN 1..27 LOOP
    v_result := v_result || SUBSTR(v_alphabet, TRUNC(DBMS_RANDOM.VALUE(1, 63)), 1);
  END LOOP;
  RETURN v_result;
END generate_ksuid;
/

--------------------------------------------------------------------------------
-- USER
--------------------------------------------------------------------------------
CREATE TABLE "user" (
  id_user              VARCHAR2(36)  NOT NULL,
  internal_identifier  VARCHAR2(36),
  type                 VARCHAR2(10)  NOT NULL,
  name                 VARCHAR2(120) NOT NULL,
  cpf                  VARCHAR2(11)  NOT NULL,
  birth_date           TIMESTAMP     NOT NULL,
  phone_number         VARCHAR2(20)  NOT NULL,
  email                VARCHAR2(255) NOT NULL,
  password             VARCHAR2(255) NOT NULL,
  milk_donated         NUMBER,

  created_at           TIMESTAMP     NOT NULL,
  created_by           VARCHAR2(36)  NOT NULL,
  updated_at           TIMESTAMP,
  updated_by           VARCHAR2(36),
  removed_at           TIMESTAMP,
  removed_by           VARCHAR2(36),

  CONSTRAINT pk_user PRIMARY KEY (id_user),
  CONSTRAINT ck_user_type CHECK (type IN ('common', 'adm', 'nurse')),
  CONSTRAINT fk_user_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_user_removed_by FOREIGN KEY (removed_by) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- DONATION POINT
--------------------------------------------------------------------------------
CREATE TABLE donation_point (
  id_donation_point VARCHAR2(36)  DEFAULT ('dpt_' || generate_ksuid()) NOT NULL,

  name               VARCHAR2(150) NOT NULL,
  description        CLOB,
  has_home           BOOLEAN       NOT NULL,
  phone_number       VARCHAR2(20),
  email              VARCHAR2(255),
  opening_hours      VARCHAR2(255),

  removed_at         TIMESTAMP,

  CONSTRAINT pk_donation_point PRIMARY KEY (id_donation_point),
  CONSTRAINT uq_donation_point_email UNIQUE (email)
);

--------------------------------------------------------------------------------
-- ADDRESS
--------------------------------------------------------------------------------
CREATE TABLE address (
  id_address        VARCHAR2(36)  DEFAULT ('adr_' || generate_ksuid()) NOT NULL,

  id_user           VARCHAR2(36),
  id_donation_point VARCHAR2(36),

  zipcode           VARCHAR2(10)  NOT NULL, -- CEP
  street            VARCHAR2(150) NOT NULL,
  "number"          VARCHAR2(10),
  city              VARCHAR2(100) NOT NULL,
  state             VARCHAR2(2)   NOT NULL,
  neighborhood      VARCHAR2(100) NOT NULL,
  complement        VARCHAR2(100),
  latitude          NUMBER(10,7),
  longitude         NUMBER(10,7),

  created_at        TIMESTAMP     NOT NULL,
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
  id_donation       VARCHAR2(36) NOT NULL,

  quantity_donated  NUMBER(10,2),
  is_active         BOOLEAN      NOT NULL,
  user_feedback     CLOB,
  score_feedback    NUMBER(2),                         -- migration 5

  created_at        TIMESTAMP    NOT NULL,
  created_by        VARCHAR2(36) NOT NULL,
  updated_at        TIMESTAMP,
  updated_by        VARCHAR2(36),
  removed_at        TIMESTAMP,
  removed_by        VARCHAR2(36),

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
  id_donation_step VARCHAR2(36) NOT NULL,

  id_donation      VARCHAR2(36) NOT NULL,
  id_address       VARCHAR2(36),                        -- migration 4

  name             VARCHAR2(30) NOT NULL,
  description      CLOB         NOT NULL,
  status           VARCHAR2(10) NOT NULL,
  set_date         TIMESTAMP,

  created_at       TIMESTAMP    NOT NULL,
  created_by       VARCHAR2(36) NOT NULL,
  updated_at       TIMESTAMP,
  updated_by       VARCHAR2(36),
  completed_at     TIMESTAMP,

  CONSTRAINT pk_donation_step PRIMARY KEY (id_donation_step),
  CONSTRAINT ck_donation_step_name CHECK (name IN (
    'Exame de sangue', 'Entregar kit de ordenha', 'Coletar leite', 'Análise de leite'
  )),
  CONSTRAINT ck_donation_step_status CHECK (status IN (
    'pending', 'review', 'done', 'warn', 'failed'
  )),
  CONSTRAINT fk_donation_step_donation FOREIGN KEY (id_donation) REFERENCES donation(id_donation),
  CONSTRAINT fk_donation_step_address FOREIGN KEY (id_address) REFERENCES address(id_address),
  CONSTRAINT fk_donation_step_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user),
  CONSTRAINT fk_donation_step_updated_by FOREIGN KEY (updated_by) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- DONATION STEP TIMELINE
--------------------------------------------------------------------------------
CREATE TABLE donation_step_timeline (
  id_donation_step_timeline VARCHAR2(36) NOT NULL,

  id_donation_step          VARCHAR2(36) NOT NULL,
  id_address                VARCHAR2(36),                -- migration 4

  description               CLOB         NOT NULL,
  status                    VARCHAR2(10) NOT NULL,
  set_date                  TIMESTAMP,

  created_at                TIMESTAMP    NOT NULL,
  created_by                VARCHAR2(36) NOT NULL,

  CONSTRAINT pk_donation_step_timeline PRIMARY KEY (id_donation_step_timeline),
  CONSTRAINT ck_donation_step_timeline_status CHECK (status IN (
    'pending', 'review', 'done', 'warn', 'failed'
  )),
  CONSTRAINT fk_timeline_donation_step FOREIGN KEY (id_donation_step) REFERENCES donation_step(id_donation_step),
  CONSTRAINT fk_timeline_address FOREIGN KEY (id_address) REFERENCES address(id_address),
  CONSTRAINT fk_timeline_created_by FOREIGN KEY (created_by) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- JOB
--------------------------------------------------------------------------------
CREATE TABLE job (
  id_job         VARCHAR2(36) NOT NULL,

  id_user        VARCHAR2(36) NOT NULL,
  id_step        VARCHAR2(36) NOT NULL,

  status         VARCHAR2(10)  NOT NULL,
  name           VARCHAR2(120) NOT NULL,
  description    CLOB          NOT NULL,
  date_set       TIMESTAMP,
  user_feedback  CLOB,

  created_at     TIMESTAMP    NOT NULL,
  created_by     VARCHAR2(36) NOT NULL,
  updated_at     TIMESTAMP,
  updated_by     VARCHAR2(36),
  removed_at     TIMESTAMP,
  removed_by     VARCHAR2(36),

  CONSTRAINT pk_job PRIMARY KEY (id_job),
  CONSTRAINT ck_job_status CHECK (status IN ('pending', 'done', 'failed')),
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
  id_consent_log VARCHAR2(36) NOT NULL,

  id_user        VARCHAR2(36) NOT NULL,

  terms_version  VARCHAR2(50)  NOT NULL,
  accepted_at    TIMESTAMP     NOT NULL,
  ip_address     VARCHAR2(45)  NOT NULL,
  user_agent     CLOB          NOT NULL,

  CONSTRAINT pk_consent_log PRIMARY KEY (id_consent_log),
  CONSTRAINT uq_consent_log_user_version UNIQUE (id_user, terms_version),
  CONSTRAINT fk_consent_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- USER BABY
--------------------------------------------------------------------------------
CREATE TABLE user_baby (
  id_user_baby VARCHAR2(36) NOT NULL,

  id_user      VARCHAR2(36) NOT NULL,

  name         VARCHAR2(120),
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
  id         VARCHAR2(36)  DEFAULT generate_uuid() NOT NULL,

  source     VARCHAR2(100) NOT NULL,
  content    CLOB          NOT NULL,
  embedding  VECTOR(384, FLOAT32) NOT NULL,             -- requer Oracle Database 23ai
  metadata   CLOB,

  created_at TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,

  CONSTRAINT pk_kb_chunks PRIMARY KEY (id),
  CONSTRAINT ck_kb_chunks_metadata_json CHECK (metadata IS JSON)
);

--------------------------------------------------------------------------------
-- CONVERSATIONS
--------------------------------------------------------------------------------
CREATE TABLE conversations (
  id               VARCHAR2(36) DEFAULT generate_uuid() NOT NULL,

  user_id          VARCHAR2(36) NOT NULL,
  started_at       TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
  last_message_at  TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
  summary          CLOB,

  CONSTRAINT pk_conversations PRIMARY KEY (id),
  CONSTRAINT fk_conversations_user FOREIGN KEY (user_id) REFERENCES "user"(id_user)
);

--------------------------------------------------------------------------------
-- MESSAGES
--------------------------------------------------------------------------------
CREATE TABLE messages (
  id              VARCHAR2(36) DEFAULT generate_uuid() NOT NULL,

  conversation_id VARCHAR2(36) NOT NULL,
  role            VARCHAR2(20) NOT NULL,
  content         CLOB         NOT NULL,
  tokens_used     NUMBER,

  created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,

  CONSTRAINT pk_messages PRIMARY KEY (id),
  CONSTRAINT ck_messages_role CHECK (role IN ('user', 'assistant', 'system')),
  CONSTRAINT fk_messages_conversation FOREIGN KEY (conversation_id)
    REFERENCES conversations(id) ON DELETE CASCADE
);

--------------------------------------------------------------------------------
-- LLM AUDIT
--------------------------------------------------------------------------------
CREATE TABLE llm_audit (
  id              VARCHAR2(36) DEFAULT generate_uuid() NOT NULL,

  user_id         VARCHAR2(36),
  conversation_id VARCHAR2(36),
  message_id      VARCHAR2(36),
  prompt_full     CLOB         NOT NULL,
  chunks_used     CLOB,
  llm_provider    VARCHAR2(30) NOT NULL,
  llm_model       VARCHAR2(50) NOT NULL,
  tokens_input    NUMBER,
  tokens_output   NUMBER,
  latency_ms      NUMBER,

  created_at      TIMESTAMP WITH TIME ZONE DEFAULT SYSTIMESTAMP NOT NULL,
  is_anonymous    BOOLEAN      DEFAULT FALSE NOT NULL,
  session_id      VARCHAR2(36),
  ip_hash         VARCHAR2(64),
  action_emitted  VARCHAR2(30),

  CONSTRAINT pk_llm_audit PRIMARY KEY (id),
  CONSTRAINT ck_llm_audit_prompt_json CHECK (prompt_full IS JSON),
  CONSTRAINT ck_llm_audit_chunks_json CHECK (chunks_used IS JSON)
);

--------------------------------------------------------------------------------
-- ALEMBIC VERSION (controle de migrations do serviço de IA)
--------------------------------------------------------------------------------
CREATE TABLE alembic_version (
  version_num VARCHAR2(32) NOT NULL,

  CONSTRAINT pk_alembic_version PRIMARY KEY (version_num)
);

--------------------------------------------------------------------------------
-- ÍNDICES
--------------------------------------------------------------------------------
CREATE INDEX ix_conversations_user_id ON conversations (user_id);
CREATE INDEX ix_conversations_last_message_at ON conversations (last_message_at);
CREATE INDEX ix_messages_conversation_id ON messages (conversation_id);
CREATE INDEX ix_messages_created_at ON messages (created_at);
CREATE INDEX ix_llm_audit_session_id ON llm_audit (session_id);

-- Índice vetorial (Oracle AI Vector Search, requer Oracle Database 23ai).
CREATE VECTOR INDEX ix_kb_chunks_embedding ON kb_chunks (embedding)
  ORGANIZATION INMEMORY NEIGHBOR GRAPH
  DISTANCE COSINE
  WITH TARGET ACCURACY 95;
