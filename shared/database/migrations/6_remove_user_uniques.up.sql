ALTER TABLE "user"
DROP CONSTRAINT IF EXISTS user_internal_identifier_key,
DROP CONSTRAINT IF EXISTS user_cpf_key,
DROP CONSTRAINT IF EXISTS user_phone_number_key,
DROP CONSTRAINT IF EXISTS user_email_key;
