DELETE FROM alembic_version
WHERE version_num = 'd4e91a7c22b0';

DROP INDEX IF EXISTS ix_llm_audit_session_id;
DROP TABLE IF EXISTS llm_audit;

DROP INDEX IF EXISTS ix_messages_created_at;
DROP INDEX IF EXISTS ix_messages_conversation_id;
DROP TABLE IF EXISTS messages;

DROP INDEX IF EXISTS ix_conversations_last_message_at;
DROP INDEX IF EXISTS ix_conversations_user_id;
DROP TABLE IF EXISTS conversations;

DROP INDEX IF EXISTS ix_kb_chunks_embedding;
DROP TABLE IF EXISTS kb_chunks;

DROP EXTENSION IF EXISTS vector;
DROP EXTENSION IF EXISTS pgcrypto;