-- Extensão necessária (embeddings são vetores e não textos)
CREATE EXTENSION IF NOT EXISTS vector;

-- ENUMS
CREATE TYPE enum_message_role AS ENUM ('user', 'assistant');

-- CONVERSATION
CREATE TABLE conversation (
  id_conversation VARCHAR(36) PRIMARY KEY,

  id_user VARCHAR(36) NOT NULL,

  title VARCHAR(255),

  created_at TIMESTAMP NOT NULL,
  last_message_at TIMESTAMP,

  CONSTRAINT fk_conversation_user FOREIGN KEY (id_user) REFERENCES "user"(id_user)
);

-- MESSAGE
CREATE TABLE message (
  id_message VARCHAR(36) PRIMARY KEY,

  id_conversation VARCHAR(36) NOT NULL,

  role enum_message_role NOT NULL,
  content TEXT NOT NULL,
  tokens_used INTEGER,

  created_at TIMESTAMP NOT NULL,

  CONSTRAINT fk_message_conversation FOREIGN KEY (id_conversation) REFERENCES conversation(id_conversation)
);

-- KNOWLEDGE BASE CHUNKS
CREATE TABLE kb_chunk (
  id_kb_chunk VARCHAR(36) PRIMARY KEY,
  
  source VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  embedding vector(384) NOT NULL,
  metadata JSONB,
  created_at TIMESTAMP NOT NULL
);

-- INDICE PARA BUSCA SEMANTICA (indice para melhorar perfomance de busca no vetor)
CREATE INDEX kb_chunk_embedding_idx
  on kb_chunk USING ivfflat (embedding vector_cosine_ops)
  WITH (lists = 100);

-- LLM AUDIT
CREATE TABLE llm_audit (
  id_llm_audit VARCHAR(36) PRIMARY KEY,
  
  id_user VARCHAR(36) NOT NULL,
  id_conversation VARCHAR(36) NOT NULL,
  id_message VARCHAR(36) NOT NULL,

  prompt_full JSONB NOT NULL,
  chunks_used JSONB,

  llm_provider VARCHAR(120) NOT NULL,
  llm_model VARCHAR(120) NOT NULL,

  tokens_input INTEGER,
  tokens_output INTEGER,
  latency_ms INTEGER,

  created_at TIMESTAMP NOT NULL
--sem FOREIGN KEY: refs lógicas pra preservar auditoria após exclusão LGPD
);
