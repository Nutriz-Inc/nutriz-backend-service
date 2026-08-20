CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS kb_chunks (
    id          UUID        NOT NULL DEFAULT gen_random_uuid(),
    source      VARCHAR(100) NOT NULL,
    content     TEXT        NOT NULL,
    embedding   VECTOR(384) NOT NULL,
    metadata    JSONB,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT kb_chunks_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS ix_kb_chunks_embedding
    ON kb_chunks USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);

CREATE TABLE IF NOT EXISTS conversations (
    id               UUID        NOT NULL DEFAULT gen_random_uuid(),
    user_id          VARCHAR(36) NOT NULL,
    started_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_message_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    summary          TEXT,
    CONSTRAINT conversations_pkey PRIMARY KEY (id),
    CONSTRAINT conversations_user_id_fkey
        FOREIGN KEY (user_id) REFERENCES "user" (id_user)
);

CREATE INDEX IF NOT EXISTS ix_conversations_user_id
    ON conversations (user_id);
CREATE INDEX IF NOT EXISTS ix_conversations_last_message_at
    ON conversations (last_message_at);

CREATE TABLE IF NOT EXISTS messages (
    id               UUID        NOT NULL DEFAULT gen_random_uuid(),
    conversation_id  UUID        NOT NULL,
    role             VARCHAR(20) NOT NULL,
    content          TEXT        NOT NULL,
    tokens_used      INTEGER,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT messages_pkey PRIMARY KEY (id),
    CONSTRAINT messages_conversation_id_fkey
        FOREIGN KEY (conversation_id) REFERENCES conversations (id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS ix_messages_conversation_id
    ON messages (conversation_id);
CREATE INDEX IF NOT EXISTS ix_messages_created_at
    ON messages (created_at);

CREATE TABLE IF NOT EXISTS llm_audit (
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
    CONSTRAINT llm_audit_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS ix_llm_audit_session_id
    ON llm_audit (session_id);

CREATE TABLE IF NOT EXISTS alembic_version (
    version_num VARCHAR(32) NOT NULL,
    CONSTRAINT alembic_version_pkc PRIMARY KEY (version_num)
);

INSERT INTO alembic_version (version_num)
    SELECT 'd4e91a7c22b0'
    WHERE NOT EXISTS (SELECT 1 FROM alembic_version);