--------------------------------------------------------------------------------
-- nutriz_indicators.sql
-- Indicadores de negócio sobre o banco de dados Nutriz (Oracle SQL).
--
-- Pré-requisitos: executar, nessa ordem, nutriz.sql (DDL) e nutriz_seed.sql
-- (massa de dados).
--
-- 8 consultas cobrindo GROUP BY, HAVING, funções de agregação e subqueries,
-- cada uma representando um indicador relevante para a operação da Nutriz.
--------------------------------------------------------------------------------

-- 1) Ranking de doadoras: volume total doado e número de doações concluídas.
--    Indicador: identifica as maiores doadoras para ações de reconhecimento/fidelização.
SELECT u.name AS doadora,
       COUNT(d.id_donation) AS qtd_doacoes,
       SUM(d.quantity_donated) AS total_doado_ml
FROM donation d
INNER JOIN "user" u ON u.id_user = d.created_by
WHERE d.quantity_donated IS NOT NULL
GROUP BY u.name
ORDER BY total_doado_ml DESC;

-- 2) Doadoras recorrentes (mais de uma doação registrada, concluída ou em andamento).
--    Indicador: mede fidelização/retenção de doadoras.
SELECT u.name AS doadora, COUNT(*) AS qtd_doacoes
FROM donation d
INNER JOIN "user" u ON u.id_user = d.created_by
GROUP BY u.name
HAVING COUNT(*) > 1
ORDER BY qtd_doacoes DESC;

-- 3) Avaliação média por posto de coleta, somente postos com pelo menos 2 doações avaliadas.
--    Indicador: qualidade percebida do atendimento por unidade, com amostra mínima confiável.
SELECT dp.name AS posto,
       COUNT(DISTINCT d.id_donation) AS qtd_doacoes_avaliadas,
       ROUND(AVG(d.score_feedback), 2) AS media_avaliacao
FROM donation d
INNER JOIN donation_step ds ON ds.id_donation = d.id_donation AND ds.name = 'Análise de leite'
INNER JOIN address a ON a.id_address = ds.id_address
INNER JOIN donation_point dp ON dp.id_donation_point = a.id_donation_point
WHERE d.score_feedback IS NOT NULL
GROUP BY dp.name
HAVING COUNT(DISTINCT d.id_donation) >= 2
ORDER BY media_avaliacao DESC;

-- 4) Volume de etapas do pipeline de doação por nome e status.
--    Indicador: mostra em quais etapas há mais gargalos (ex.: muitas etapas "pending").
SELECT name AS etapa, status, COUNT(*) AS qtd
FROM donation_step
GROUP BY name, status
ORDER BY etapa, status;

-- 5) Produtividade da equipe: enfermeiras com 2 ou mais visitas (jobs) concluídas.
--    Indicador: apoia decisões de alocação de equipe para visitas domiciliares.
SELECT u.name AS enfermeira, COUNT(*) AS jobs_concluidos
FROM job j
INNER JOIN "user" u ON u.id_user = j.id_user
WHERE j.status = 'done'
GROUP BY u.name
HAVING COUNT(*) >= 2
ORDER BY jobs_concluidos DESC;

-- 6) Doadoras cujo volume total doado está acima da média geral (subquery escalar).
--    Indicador: público-alvo para reconhecimento/depoimentos de maiores doadoras.
SELECT name AS doadora, milk_donated
FROM "user"
WHERE milk_donated > (
  SELECT AVG(milk_donated) FROM "user" WHERE milk_donated IS NOT NULL
)
ORDER BY milk_donated DESC;

-- 7) Postos de coleta com pelo menos uma doação avaliada com nota máxima (subquery correlacionada com EXISTS).
--    Indicador: unidades com casos de excelência no atendimento, úteis para benchmarking interno.
SELECT dp.name AS posto
FROM donation_point dp
WHERE EXISTS (
  SELECT 1
  FROM donation_step ds
  INNER JOIN address a ON a.id_address = ds.id_address
  INNER JOIN donation dn ON dn.id_donation = ds.id_donation
  WHERE a.id_donation_point = dp.id_donation_point
    AND ds.name = 'Análise de leite'
    AND dn.score_feedback = 5
)
ORDER BY dp.name;

-- 8) Doadoras cadastradas que ainda não iniciaram nenhuma doação (subquery com NOT EXISTS).
--    Indicador: funil de conversão de cadastro para primeira doação, base para campanhas de ativação.
SELECT u.name AS doadora, u.email
FROM "user" u
WHERE u.type = 'common'
  AND NOT EXISTS (
    SELECT 1 FROM donation d WHERE d.created_by = u.id_user
  )
ORDER BY u.name;
