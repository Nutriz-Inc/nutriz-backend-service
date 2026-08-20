--------------------------------------------------------------------------------
-- nutriz_joins.sql
-- Consultas com JOIN sobre o banco de dados Nutriz (Oracle SQL).
--
-- Pré-requisitos: executar, nessa ordem, nutriz.sql (DDL) e nutriz_seed.sql
-- (massa de dados).
--
-- 6 consultas: 2 com INNER JOIN, 2 com LEFT JOIN e 2 com RIGHT JOIN.
-- RIGHT JOIN é usado nos casos em que a tabela "preservada" (que deve
-- aparecer inteira, mesmo sem correspondência) fica à direita da junção;
-- nesses dois casos ela é logicamente equivalente a inverter as tabelas e
-- usar LEFT JOIN, mas está aqui para demonstrar a sintaxe, suportada
-- nativamente pelo Oracle.
--------------------------------------------------------------------------------

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

-- 6) RIGHT JOIN: todas as etapas de doação e o job (visita) associado, mesmo etapas sem job vinculado
--    (ex.: "Exame de sangue" e "Análise de leite", que não geram visita domiciliar).
SELECT ds.id_donation_step, ds.name AS etapa, ds.status, j.name AS job_associado
FROM job j
RIGHT JOIN donation_step ds ON ds.id_donation_step = j.id_step
ORDER BY ds.id_donation_step;
