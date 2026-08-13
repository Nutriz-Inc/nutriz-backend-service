ALTER TABLE "user"
ADD CONSTRAINT user_internal_identifier_key UNIQUE (internal_identifier),
ADD CONSTRAINT user_cpf_key UNIQUE (cpf),
ADD CONSTRAINT user_phone_number_key UNIQUE (phone_number),
ADD CONSTRAINT user_email_key UNIQUE (email);
