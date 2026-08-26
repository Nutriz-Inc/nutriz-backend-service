# Indicadores de Negócio — Nutriz

9 indicadores construídos a partir das consultas já desenvolvidas em
[`nutriz_queries.sql`](nutriz_queries.sql), [`nutriz_joins.sql`](nutriz_joins.sql)
e [`nutriz_indicators.sql`](nutriz_indicators.sql). Alguns são a própria consulta
original; outros a envolvem em uma agregação simples para chegar a um número
executivo — em ambos os casos, a fonte é indicada em cada item.

Resultados apurados manualmente sobre a massa de dados de
[`nutriz_seed.sql`](nutriz_seed.sql). Execute `nutriz.sql` (ou
`nutriz-postgres.sql`) → o seed correspondente → a consulta, em um banco real,
para reproduzir os mesmos números.

---

## 1. Ranking de doadoras por volume total doado

**Fonte:** `nutriz_indicators.sql`, consulta 1.

**Objetivo do indicador**
Identificar as doadoras que mais contribuíram em volume de leite (ml), para
priorizar ações de reconhecimento e fidelização.

**Consulta SQL utilizada**
```sql
SELECT u.name AS doadora,
       COUNT(d.id_donation) AS qtd_doacoes,
       SUM(d.quantity_donated) AS total_doado_ml
FROM donation d
INNER JOIN "user" u ON u.id_user = d.created_by
WHERE d.quantity_donated IS NOT NULL
GROUP BY u.name
ORDER BY total_doado_ml DESC;
```

**Resultado obtido**

| DOADORA | QTD_DOACOES | TOTAL_DOADO_ML |
|---|---|---|
| Beatriz Alves | 1 | 1200.00 |
| Gabriela Teixeira | 1 | 900.00 |
| Maria Silva | 1 | 850.50 |
| Joana Souza | 1 | 600.00 |
| Paula Martins | 1 | 450.25 |

**Benefício para o negócio**
Permite programas de reconhecimento (selos, depoimentos, brindes) direcionados
às maiores doadoras, incentivando a manutenção do volume de doações e servindo
de prova social para atrair novas doadoras.

---

## 2. Doadoras recorrentes (mais de uma doação)

**Fonte:** `nutriz_indicators.sql`, consulta 2.

**Objetivo do indicador**
Medir a taxa de retenção: quantas doadoras voltam a doar mais de uma vez.

**Consulta SQL utilizada**
```sql
SELECT u.name AS doadora, COUNT(*) AS qtd_doacoes
FROM donation d
INNER JOIN "user" u ON u.id_user = d.created_by
GROUP BY u.name
HAVING COUNT(*) > 1
ORDER BY qtd_doacoes DESC;
```

**Resultado obtido**

| DOADORA | QTD_DOACOES |
|---|---|
| Maria Silva | 2 |

**Benefício para o negócio**
Apenas 1 de 10 doadoras cadastradas já iniciou uma segunda doação — sinal de
que o funil de retenção pós-primeira-doação precisa de reforço (ex.: lembretes
automáticos, campanhas de reengajamento após alguns meses).

---

## 3. Cobertura de perfil do bebê entre as doadoras

**Fonte:** `nutriz_joins.sql`, consulta 3 (`"user" LEFT JOIN user_baby`),
agregada para obter o total comparável.

**Objetivo do indicador**
Medir quantas doadoras completaram o cadastro do bebê — dado usado para
elegibilidade e para o app sugerir conteúdo relevante à fase da amamentação.

**Consulta SQL utilizada**
```sql
SELECT
  COUNT(DISTINCT CASE WHEN u.type = 'common' THEN u.id_user END) AS total_doadoras,
  COUNT(DISTINCT CASE WHEN u.type = 'common' AND b.id_user IS NOT NULL
                  THEN u.id_user END) AS doadoras_com_bebe
FROM "user" u
LEFT JOIN user_baby b ON b.id_user = u.id_user;
```

**Resultado obtido**

| TOTAL_DOADORAS | DOADORAS_COM_BEBE |
|---|---|
| 10 | 7 (70%) |

Sem bebê cadastrado: Amanda Rocha, Rafaela Pinto e Débora Nunes.

**Benefício para o negócio**
70% de completude é um número acionável: as 3 doadoras sem bebê cadastrado
viram uma lista direta para um lembrete de "complete seu perfil" dentro do
app, reduzindo cadastros incompletos.

---

## 4. Volume de etapas do pipeline de doação por status

**Fonte:** `nutriz_indicators.sql`, consulta 4.

**Objetivo do indicador**
Enxergar em que etapa do processo (exame, kit, coleta, análise) as doações
estão concentradas, revelando gargalos operacionais.

**Consulta SQL utilizada**
```sql
SELECT name AS etapa, status, COUNT(*) AS qtd
FROM donation_step
GROUP BY name, status
ORDER BY etapa, status;
```

**Resultado obtido**

| ETAPA | STATUS | QTD |
|---|---|---|
| Análise de leite | done | 5 |
| Coletar leite | done | 5 |
| Entregar kit de ordenha | done | 5 |
| Entregar kit de ordenha | pending | 5 |
| Exame de sangue | done | 10 |

**Benefício para o negócio**
As 5 doações em andamento estão todas paradas na mesma etapa (entrega do kit
pendente), enquanto o exame de sangue já foi concluído para todas. Isso
identifica a entrega do kit como o gargalo atual do pipeline.

---

## 5. Cobertura de visitas (job) por etapa do pipeline

**Fonte:** `nutriz_joins.sql`, consulta 6 (`job RIGHT JOIN donation_step`),
agregada.

**Objetivo do indicador**
Medir que fração das etapas do pipeline gera uma tarefa de visita domiciliar
(job) para a equipe, versus etapas puramente laboratoriais ou ainda sem
tarefa criada.

**Consulta SQL utilizada**
```sql
SELECT
  COUNT(*) AS total_etapas,
  COUNT(CASE WHEN j.id_job IS NOT NULL THEN 1 END) AS etapas_com_job,
  COUNT(CASE WHEN j.id_job IS NULL THEN 1 END) AS etapas_sem_job
FROM job j
RIGHT JOIN donation_step ds ON ds.id_donation_step = j.id_step;
```

**Resultado obtido**

| TOTAL_ETAPAS | ETAPAS_COM_JOB | ETAPAS_SEM_JOB |
|---|---|---|
| 30 | 10 (33%) | 20 (67%) |

**Benefício para o negócio**
Confirma que "Exame de sangue" e "Análise de leite" nunca geram visita (são
etapas de laboratório), enquanto as 5 etapas de kit ainda pendentes das
doações em andamento **também não têm job criado** — um ponto operacional a
corrigir junto do gargalo identificado no indicador 4.

---

## 6. Produtividade da equipe: enfermeiras com mais visitas concluídas

**Fonte:** `nutriz_indicators.sql`, consulta 5.

**Objetivo do indicador**
Medir a carga de trabalho concluída por enfermeira (entregas de kit e coletas
domiciliares), apoiando decisões de alocação de equipe.

**Consulta SQL utilizada**
```sql
SELECT u.name AS enfermeira, COUNT(*) AS jobs_concluidos
FROM job j
INNER JOIN "user" u ON u.id_user = j.id_user
WHERE j.status = 'done'
GROUP BY u.name
HAVING COUNT(*) >= 2
ORDER BY jobs_concluidos DESC;
```

**Resultado obtido**

| ENFERMEIRA | JOBS_CONCLUIDOS |
|---|---|
| Fernanda Duarte | 4 |
| Juliana Ramos | 4 |
| Patrícia Gomes | 2 |

**Benefício para o negócio**
Mostra que a carga de visitas está concentrada em duas enfermeiras (Fernanda
e Juliana), enquanto Patrícia tem menos da metade do volume — indicativo para
rebalancear a distribuição de doadoras entre a equipe.

---

## 7. Cobertura de coleta domiciliar entre os postos de doação

**Fonte:** `nutriz_queries.sql`, consulta 6, agregada.

**Objetivo do indicador**
Medir quantos dos postos parceiros oferecem coleta domiciliar — um
diferencial que reduz o atrito da doadora e tende a aumentar a retenção
(indicador 2).

**Consulta SQL utilizada**
```sql
SELECT
  COUNT(*) AS total_postos,
  COUNT(CASE WHEN has_home = TRUE THEN 1 END) AS postos_com_coleta_domiciliar
FROM donation_point;
```

**Resultado obtido**

| TOTAL_POSTOS | POSTOS_COM_COLETA_DOMICILIAR |
|---|---|
| 5 | 1 (20%) |

Único posto com coleta em domicílio: Banco de Leite Humano do Hospital Ipiranga.

**Benefício para o negócio**
Só 20% da rede oferece coleta domiciliar hoje. Como esse benefício reduz
atrito, é um argumento de negociação concreto para expandir o serviço junto
aos demais parceiros (Amparo Maternal, São Luiz Star, Anália Franco e Santa
Casa).

---

## 8. Doadoras cadastradas que ainda não iniciaram nenhuma doação

**Fonte:** `nutriz_indicators.sql`, consulta 8.

**Objetivo do indicador**
Medir o funil de conversão entre cadastro no app e a primeira doação
efetivamente iniciada.

**Consulta SQL utilizada**
```sql
SELECT u.name AS doadora, u.email
FROM "user" u
WHERE u.type = 'common'
  AND NOT EXISTS (
    SELECT 1 FROM donation d WHERE d.created_by = u.id_user
  )
ORDER BY u.name;
```

**Resultado obtido**

| DOADORA | EMAIL |
|---|---|
| Débora Nunes | debora.nunes@gmail.com |

**Benefício para o negócio**
De 10 usuárias comuns cadastradas, apenas 1 nunca iniciou uma doação — 90% de
conversão. Essa é exatamente a lista que uma campanha de ativação
(notificação, e-mail, ligação da equipe) deve acionar primeiro.

---

## 9. Concentração do stack de IA em produção

**Fonte:** `nutriz_queries.sql`, consulta 8, envolvida em `COUNT(*)`.

**Objetivo do indicador**
Verificar quantos provedores/modelos de IA distintos estão de fato em uso
pelo assistente virtual — relevante para avaliar risco de dependência de um
único fornecedor (vendor lock-in) e ausência de fallback.

**Consulta SQL utilizada**
```sql
SELECT COUNT(*) AS combinacoes_distintas
FROM (
  SELECT DISTINCT llm_provider, llm_model
  FROM llm_audit
);
```

**Resultado obtido**

| COMBINACOES_DISTINTAS |
|---|
| 1 |

As 10 interações auditadas usam sempre `anthropic` / `claude-sonnet-5`.

**Benefício para o negócio**
Confirma 100% de dependência de um único provedor/modelo. É um indicador de
risco operacional: vale a pena avaliar um modelo de fallback (outro provedor
ou versão) para continuidade do assistente em caso de indisponibilidade ou
descontinuação do modelo atual.
