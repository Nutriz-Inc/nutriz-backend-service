-- ENUMS
CREATE TYPE enum_message_role AS ENUM ('user', 'assistant');

-- CONVERSATION
CREATE TABLE conversation (
  id_conversation VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36) NOT NULL,

  title VARCHAR(255) NOT NULL,

  created_at TIMESTAMP NOT NULL,
  last_message_at TIMESTAMP NOT NULL,

  CONSTRAINT fk_conversation_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);

-- MESSAGE
CREATE TABLE message (
  id_message VARCHAR(36) PRIMARY KEY,

  id_conversation VARCHAR(36) NOT NULL,

  role enum_message_role NOT NULL,
  content TEXT NOT NULL,
  tokens_used NUMERIC,

  created_at TIMESTAMP NOT NULL,

  CONSTRAINT fk_message_conversation FOREIGN KEY (id_conversation) REFERENCES conversation(id_conversation)
);

-- KNOWLEDGE BASE CHUNKS
CREATE TABLE kb_chunk (
  id_kb_chunk VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36) NOT NULL,

  source VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  embedding TEXT NOT NULL,
  metadata TEXT,

  created_at TIMESTAMP NOT NULL,

  CONSTRAINT fk_kb_chunk_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);

-- LLM AUDIT
CREATE TABLE llm_audit (
  id_llm_audit VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36) NOT NULL,
  id_message VARCHAR(36) NOT NULL,

  prompt_full TEXT NOT NULL,
  chunks_used TEXT,

  llm_provider VARCHAR(120) NOT NULL,
  llm_model VARCHAR(120) NOT NULL,

  tokens_input NUMERIC,
  tokens_output NUMERIC,
  latency_ms NUMERIC,

  created_at TIMESTAMP NOT NULL,

  CONSTRAINT fk_llm_audit_user FOREIGN KEY (id_user) REFERENCES "user"(id_user),

  CONSTRAINT fk_llm_audit_message FOREIGN KEY (id_message) REFERENCES message(id_message)
);