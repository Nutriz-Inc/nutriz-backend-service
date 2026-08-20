# Indicadores de Negócio — Nutriz

Os resultados abaixo foram apurados manualmente a partir da massa de dados de
[`nutriz_seed.sql`](nutriz_seed.sql), aplicando as consultas de
[`nutriz_indicators.sql`](nutriz_indicators.sql). Executando `nutriz.sql` →
`nutriz_seed.sql` → a consulta correspondente em uma instância Oracle, os
mesmos resultados devem ser reproduzidos.

---

## 1. Ranking de doadoras por volume total doado

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
Apenas 1 de 10 doadoras cadastradas já iniciou uma segunda doação — um sinal
de que o funil de retenção pós-primeira-doação precisa de reforço (ex.:
lembretes automáticos, campanhas de reengajamento após alguns meses).

---

## 3. Avaliação média por posto de coleta (amostra mínima de 2 doações)

**Objetivo do indicador**
Comparar a qualidade percebida do atendimento entre os postos de coleta,
descartando postos com amostra insuficiente (< 2 avaliações) para evitar
conclusões enviesadas por um único caso.

**Consulta SQL utilizada**
```sql
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
```

**Resultado obtido**

| POSTO | QTD_DOACOES_AVALIADAS | MEDIA_AVALIACAO |
|---|---|---|
| Posto de Coleta de Leite Humano Amparo Maternal | 2 | 4.00 |

**Benefício para o negócio**
Com a massa de dados atual, apenas o posto Amparo Maternal atingiu a amostra
mínima (2 doações concluídas e avaliadas: notas 5 e 3). Isso já aponta uma
nota "puxada para baixo" por uma avaliação de 3 — um gatilho para a equipe
investigar o que aconteceu naquele atendimento específico antes que vire
padrão.

---

## 4. Volume de etapas do pipeline de doação por status

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
As 5 doações em andamento estão todas paradas na mesma etapa
("Entregar kit de ordenha" pendente), enquanto o exame de sangue já foi
concluído para todas. Isso identifica a entrega do kit como o gargalo atual
do pipeline, orientando a equipe operacional a priorizar esse agendamento
para não perder o engajamento inicial da doadora.

---

## 5. Produtividade da equipe: enfermeiras com mais visitas concluídas

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
Mostra que a carga de visitas domiciliares está concentrada em duas
enfermeiras (Fernanda e Juliana), enquanto Patrícia tem menos da metade do
volume — um indicativo para rebalancear a distribuição de doadoras entre a
equipe ou investigar se há diferença na região/dificuldade de cada carteira.

---

## 6. Doadoras com volume acima da média geral

**Objetivo do indicador**
Identificar doadoras "acima da média" para ações direcionadas (ex.: convite
para virar embaixadora, entrevista para conteúdo institucional).

**Consulta SQL utilizada**
```sql
SELECT name AS doadora, milk_donated
FROM "user"
WHERE milk_donated > (
  SELECT AVG(milk_donated) FROM "user" WHERE milk_donated IS NOT NULL
)
ORDER BY milk_donated DESC;
```

**Resultado obtido**

Média geral (subquery): (1450 + 600 + 1200 + 450 + 900) / 5 = **920 ml**

| DOADORA | MILK_DONATED |
|---|---|
| Maria Silva | 1450 |
| Beatriz Alves | 1200 |

**Benefício para o negócio**
Apenas 2 das 5 doadoras com histórico registrado superam a média — dá à
equipe de marketing uma lista curta e objetiva de perfis de destaque para
campanhas, sem depender de análise manual da base completa.

---

## 7. Doadoras cadastradas que ainda não iniciaram nenhuma doação

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
De 10 usuárias comuns cadastradas, apenas 1 nunca iniciou uma doação —
90% de conversão. Essa é exatamente a lista que uma campanha de ativação
(notificação, e-mail, ligação da equipe) deve acionar primeiro.

---

## 8. Postos de coleta com pelo menos uma doação nota máxima

**Objetivo do indicador**
Identificar unidades com casos de excelência (nota 5) no atendimento, para
uso como referência interna de boas práticas.

**Consulta SQL utilizada**
```sql
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
```

**Resultado obtido**

| POSTO |
|---|
| Banco de Leite Humano do Hospital Ipiranga |
| Banco de Leite Humano Rede Dor São Luiz - Unidade Anália Franco |
| Posto de Coleta de Leite Humano Amparo Maternal |

**Benefício para o negócio**
3 dos 5 postos já têm ao menos um atendimento nota máxima registrado — bons
candidatos para mapear o que estão fazendo certo (tempo de resposta, contato
com a doadora) e replicar nos postos que ainda não atingiram essa marca
(São Luiz Star e Santa Casa).
