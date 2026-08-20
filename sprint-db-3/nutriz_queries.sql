--------------------------------------------------------------------------------
-- nutriz_queries.sql
-- Consultas de exemplo sobre o banco de dados Nutriz (Oracle SQL).
--
-- Pré-requisitos: executar, nessa ordem, nutriz.sql (DDL) e nutriz_seed.sql
-- (massa de dados).
--
-- 8 consultas cobrindo SELECT, WHERE, ORDER BY e DISTINCT.
--------------------------------------------------------------------------------

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

-- 6) Pontos de doação que oferecem coleta de leite em domicílio.
SELECT id_donation_point, name, phone_number
FROM donation_point
WHERE has_home = TRUE
ORDER BY name;

-- 7) Respostas mais recentes do assistente de IA às doadoras.
SELECT id, conversation_id, content, created_at
FROM messages
WHERE role = 'assistant'
ORDER BY created_at DESC;

-- 8) Combinações distintas de provedor/modelo de IA usadas nas auditorias.
SELECT DISTINCT llm_provider, llm_model
FROM llm_audit
ORDER BY llm_provider, llm_model;
