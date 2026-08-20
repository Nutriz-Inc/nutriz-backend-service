--------------------------------------------------------------------------------
-- nutriz_seed.sql
-- Massa de dados inicial (carga) para o banco de dados Nutriz.
--
-- Pré-requisito: executar nutriz.sql antes deste script (cria tabelas e
-- constraints no Oracle).
--
-- Cenário de negócio: Nutriz conecta mães lactantes que desejam doar leite
-- humano excedente a bancos/postos de coleta de leite humano em São Paulo.
-- Cada doação passa por um pipeline de etapas (exame de sangue, entrega do
-- kit de ordenha, coleta do leite e análise), acompanhado por enfermeiras
-- (jobs), com consentimento (LGPD) e um assistente de IA que responde
-- dúvidas das doadoras com base em uma base de conhecimento (kb_chunks).
--
-- Distribuição de registros (total >= 100, integridade referencial garantida
-- pela ordem de inserção, que respeita as dependências de FK):
--   "user"                 15
--   donation_point          5
--   address                15
--   donation                10
--   donation_step           30
--   donation_step_timeline  30
--   job                     10
--   consent_log             15
--   user_baby                8
--   kb_chunks                6
--   conversations            8
--   messages                20
--   llm_audit               10
--   alembic_version          1
--   TOTAL                  183
--------------------------------------------------------------------------------

--------------------------------------------------------------------------------
-- USER (2 adm, 3 nurse, 10 common)
--------------------------------------------------------------------------------
INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_adm_carla', 'adm', 'Carla Menezes', '12345678901', TIMESTAMP '1985-04-12 00:00:00', '11987650001', 'carla.menezes@nutriz.org', '$2b$12$placeholderHash0000000000000000000001', TIMESTAMP '2025-01-05 09:00:00', 'system');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_adm_bruno', 'adm', 'Bruno Ferreira', '12345678902', TIMESTAMP '1988-09-23 00:00:00', '11987650002', 'bruno.ferreira@nutriz.org', '$2b$12$placeholderHash0000000000000000000002', TIMESTAMP '2025-01-05 09:15:00', 'usr_adm_carla');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_nurse_juliana', 'nurse', 'Juliana Ramos', '12345678903', TIMESTAMP '1990-01-30 00:00:00', '11987650003', 'juliana.ramos@nutriz.org', '$2b$12$placeholderHash0000000000000000000003', TIMESTAMP '2025-01-10 10:00:00', 'usr_adm_carla');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_nurse_patricia', 'nurse', 'Patrícia Gomes', '12345678904', TIMESTAMP '1992-06-18 00:00:00', '11987650004', 'patricia.gomes@nutriz.org', '$2b$12$placeholderHash0000000000000000000004', TIMESTAMP '2025-01-10 10:15:00', 'usr_adm_carla');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_nurse_fernanda', 'nurse', 'Fernanda Duarte', '12345678905', TIMESTAMP '1991-11-05 00:00:00', '11987650005', 'fernanda.duarte@nutriz.org', '$2b$12$placeholderHash0000000000000000000005', TIMESTAMP '2025-01-12 11:00:00', 'usr_adm_bruno');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_maria_silva', 'common', 'Maria Silva', '12345678906', TIMESTAMP '1994-02-14 00:00:00', '11987650006', 'maria.silva@gmail.com', '$2b$12$placeholderHash0000000000000000000006', 1450, TIMESTAMP '2025-02-01 08:30:00', 'usr_maria_silva');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_joana_souza', 'common', 'Joana Souza', '12345678907', TIMESTAMP '1996-07-22 00:00:00', '11987650007', 'joana.souza@gmail.com', '$2b$12$placeholderHash0000000000000000000007', 600, TIMESTAMP '2025-02-03 09:10:00', 'usr_joana_souza');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_camila_lima', 'common', 'Camila Lima', '12345678908', TIMESTAMP '1993-03-09 00:00:00', '11987650008', 'camila.lima@gmail.com', '$2b$12$placeholderHash0000000000000000000008', TIMESTAMP '2025-02-10 14:20:00', 'usr_camila_lima');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_beatriz_alves', 'common', 'Beatriz Alves', '12345678909', TIMESTAMP '1990-12-01 00:00:00', '11987650009', 'beatriz.alves@gmail.com', '$2b$12$placeholderHash0000000000000000000009', 1200, TIMESTAMP '2025-02-15 16:00:00', 'usr_beatriz_alves');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_larissa_costa', 'common', 'Larissa Costa', '12345678910', TIMESTAMP '1995-05-17 00:00:00', '11987650010', 'larissa.costa@gmail.com', '$2b$12$placeholderHash0000000000000000000010', TIMESTAMP '2025-03-01 08:00:00', 'usr_larissa_costa');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_paula_martins', 'common', 'Paula Martins', '12345678911', TIMESTAMP '1989-08-25 00:00:00', '11987650011', 'paula.martins@gmail.com', '$2b$12$placeholderHash0000000000000000000011', 450, TIMESTAMP '2025-03-05 09:45:00', 'usr_paula_martins');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_amanda_rocha', 'common', 'Amanda Rocha', '12345678912', TIMESTAMP '1997-10-03 00:00:00', '11987650012', 'amanda.rocha@gmail.com', '$2b$12$placeholderHash0000000000000000000012', TIMESTAMP '2025-03-10 11:30:00', 'usr_amanda_rocha');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, milk_donated, created_at, created_by)
VALUES ('usr_gabriela_teixeira', 'common', 'Gabriela Teixeira', '12345678913', TIMESTAMP '1992-01-19 00:00:00', '11987650013', 'gabriela.teixeira@gmail.com', '$2b$12$placeholderHash0000000000000000000013', 900, TIMESTAMP '2025-03-15 13:00:00', 'usr_gabriela_teixeira');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_rafaela_pinto', 'common', 'Rafaela Pinto', '12345678914', TIMESTAMP '1994-09-11 00:00:00', '11987650014', 'rafaela.pinto@gmail.com', '$2b$12$placeholderHash0000000000000000000014', TIMESTAMP '2025-03-20 15:15:00', 'usr_rafaela_pinto');

INSERT INTO "user" (id_user, type, name, cpf, birth_date, phone_number, email, password, created_at, created_by)
VALUES ('usr_debora_nunes', 'common', 'Débora Nunes', '12345678915', TIMESTAMP '1991-04-28 00:00:00', '11987650015', 'debora.nunes@gmail.com', '$2b$12$placeholderHash0000000000000000000015', TIMESTAMP '2025-04-02 10:00:00', 'usr_debora_nunes');

--------------------------------------------------------------------------------
-- DONATION POINT (bancos/postos de leite humano reais de São Paulo)
--------------------------------------------------------------------------------
INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_amparo_maternal', 'Posto de Coleta de Leite Humano Amparo Maternal', 'Posto de coleta de leite humano localizado no Amparo Maternal, Vila Clementino.', FALSE, '1150898277', 'ouvidoria@ham.spdmpais.org.br', 'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 17:00');

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_sao_luiz_star', 'Banco de Leite Humano Maternidade São Luiz Star', 'Banco de leite humano localizado na Maternidade São Luiz Star, Vila Olímpia.', FALSE, '1121211349', 'consultoriamaternidade@maternidadesaoluiz.com.br', 'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00');

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_analia_franco', 'Banco de Leite Humano Rede Dor São Luiz - Unidade Anália Franco', 'Banco de leite humano localizado no Hospital e Maternidade São Luiz, Unidade Anália Franco, Tatuapé.', FALSE, '1133861315', 'faleconosco.sadt@saoluiz.com.br', 'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00');

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_hosp_ipiranga', 'Banco de Leite Humano do Hospital Ipiranga', 'Banco de leite humano localizado no Hospital Ipiranga, 8º andar.', TRUE, '1120677866', 'hi.cemed@gmail.com', NULL);

INSERT INTO donation_point (id_donation_point, name, description, has_home, phone_number, email, opening_hours)
VALUES ('dpt_santa_casa', 'Banco de Leite Humano da Santa Casa de São Paulo', 'Banco de leite humano localizado no Hospital Central da Santa Casa de São Paulo, Vila Buarque.', FALSE, '1121767390', 'ouvidoria@santacasasp.org.br', 'Seg-Sex, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00');

--------------------------------------------------------------------------------
-- ADDRESS (5 dos postos de doação + 10 residências das doadoras)
--------------------------------------------------------------------------------
INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_dpt_amparo', 'dpt_amparo_maternal', '04040033', 'Rua Loefgren', '101', 'São Paulo', 'SP', 'Vila Clementino', -23.600032, -46.643398, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_dpt_saoluiz', 'dpt_sao_luiz_star', '04552050', 'Rua Helena', '29', 'São Paulo', 'SP', 'Vila Olímpia', -23.590756, -46.673653, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_dpt_analia', 'dpt_analia_franco', '03313000', 'Rua Francisco Marengo', '1312', 'São Paulo', 'SP', 'Tatuapé', -23.548886, -46.558248, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, complement, latitude, longitude, created_at)
VALUES ('adr_dpt_ipiranga', 'dpt_hosp_ipiranga', '04262000', 'Avenida Nazaré', '28', 'São Paulo', 'SP', 'Ipiranga', '8º andar', -23.584210, -46.611595, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_donation_point, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_dpt_santacasa', 'dpt_santa_casa', '01221020', 'Rua Dr. Cesário Mota Júnior', '112', 'São Paulo', 'SP', 'Vila Buarque', -23.542626, -46.650065, TIMESTAMP '2025-01-01 08:00:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_maria', 'usr_maria_silva', '05014010', 'Rua Turiassu', '500', 'São Paulo', 'SP', 'Perdizes', -23.535000, -46.682200, TIMESTAMP '2025-02-01 08:35:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_joana', 'usr_joana_souza', '04101000', 'Rua Vergueiro', '1200', 'São Paulo', 'SP', 'Vila Mariana', -23.589000, -46.633900, TIMESTAMP '2025-02-03 09:15:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_camila', 'usr_camila_lima', '03164000', 'Rua Coelho Barradas', '75', 'São Paulo', 'SP', 'Tatuapé', -23.540000, -46.575000, TIMESTAMP '2025-02-10 14:25:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_beatriz', 'usr_beatriz_alves', '04275000', 'Avenida do Cursino', '320', 'São Paulo', 'SP', 'Ipiranga', -23.592000, -46.610000, TIMESTAMP '2025-02-15 16:05:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_larissa', 'usr_larissa_costa', '01223000', 'Rua Barra Funda', '410', 'São Paulo', 'SP', 'Barra Funda', -23.527000, -46.662000, TIMESTAMP '2025-03-01 08:05:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_paula', 'usr_paula_martins', '05407000', 'Rua Cardeal Arcoverde', '1500', 'São Paulo', 'SP', 'Pinheiros', -23.567000, -46.691000, TIMESTAMP '2025-03-05 09:50:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_amanda', 'usr_amanda_rocha', '04547000', 'Rua Funchal', '160', 'São Paulo', 'SP', 'Vila Olímpia', -23.595000, -46.689000, TIMESTAMP '2025-03-10 11:35:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_gabriela', 'usr_gabriela_teixeira', '03310000', 'Rua Serra de Botucatu', '850', 'São Paulo', 'SP', 'Tatuapé', -23.541000, -46.568000, TIMESTAMP '2025-03-15 13:05:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_rafaela', 'usr_rafaela_pinto', '04261000', 'Avenida Nazaré', '900', 'São Paulo', 'SP', 'Ipiranga', -23.586000, -46.610000, TIMESTAMP '2025-03-20 15:20:00');

INSERT INTO address (id_address, id_user, zipcode, street, "number", city, state, neighborhood, latitude, longitude, created_at)
VALUES ('adr_usr_debora', 'usr_debora_nunes', '01310000', 'Avenida Paulista', '2000', 'São Paulo', 'SP', 'Bela Vista', -23.561000, -46.656000, TIMESTAMP '2025-04-02 10:05:00');

--------------------------------------------------------------------------------
-- DONATION (10: 5 concluídas, 5 em andamento)
--------------------------------------------------------------------------------
INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('dnt_maria_01', 850.50, FALSE, 'Processo tranquilo, equipe muito atenciosa.', 5, TIMESTAMP '2025-11-10 09:00:00', 'usr_maria_silva', TIMESTAMP '2025-11-20 17:00:00', 'usr_nurse_juliana');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('dnt_maria_02', TRUE, TIMESTAMP '2026-07-05 09:00:00', 'usr_maria_silva');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('dnt_joana_01', 600.00, FALSE, 'Gostei muito de poder ajudar bebês prematuros.', 4, TIMESTAMP '2025-11-15 10:00:00', 'usr_joana_souza', TIMESTAMP '2025-11-25 16:30:00', 'usr_nurse_juliana');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('dnt_camila_01', TRUE, TIMESTAMP '2026-06-20 11:00:00', 'usr_camila_lima');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('dnt_beatriz_01', 1200.00, FALSE, 'Coleta em domicílio facilitou muito.', 5, TIMESTAMP '2025-12-01 08:30:00', 'usr_beatriz_alves', TIMESTAMP '2025-12-10 18:00:00', 'usr_nurse_patricia');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('dnt_larissa_01', TRUE, TIMESTAMP '2026-07-10 09:30:00', 'usr_larissa_costa');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('dnt_paula_01', 450.25, FALSE, 'Poderia ser um pouco mais rápido, mas ok.', 3, TIMESTAMP '2025-12-05 13:00:00', 'usr_paula_martins', TIMESTAMP '2025-12-18 12:00:00', 'usr_nurse_fernanda');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('dnt_amanda_01', TRUE, TIMESTAMP '2026-07-15 14:00:00', 'usr_amanda_rocha');

INSERT INTO donation (id_donation, quantity_donated, is_active, user_feedback, score_feedback, created_at, created_by, updated_at, updated_by)
VALUES ('dnt_gabriela_01', 900.00, FALSE, 'Excelente atendimento do início ao fim.', 5, TIMESTAMP '2025-12-10 09:00:00', 'usr_gabriela_teixeira', TIMESTAMP '2025-12-22 17:30:00', 'usr_nurse_fernanda');

INSERT INTO donation (id_donation, is_active, created_at, created_by)
VALUES ('dnt_rafaela_01', TRUE, TIMESTAMP '2026-07-20 10:00:00', 'usr_rafaela_pinto');

--------------------------------------------------------------------------------
-- DONATION STEP (30: 4 etapas para cada doação concluída, 2 para as em andamento)
--------------------------------------------------------------------------------
-- dnt_maria_01 (posto: Amparo Maternal / residência: adr_usr_maria)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_maria01_exame', 'dnt_maria_01', 'adr_dpt_amparo', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-11-11 09:00:00', TIMESTAMP '2025-11-10 09:10:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-11 10:00:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-11 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_maria01_kit', 'dnt_maria_01', 'adr_usr_maria', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-11-13 10:00:00', TIMESTAMP '2025-11-11 10:05:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-13 10:30:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-13 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_maria01_coleta', 'dnt_maria_01', 'adr_usr_maria', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-11-18 09:00:00', TIMESTAMP '2025-11-13 10:35:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-18 09:40:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-18 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_maria01_analise', 'dnt_maria_01', 'adr_dpt_amparo', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-11-20 15:00:00', TIMESTAMP '2025-11-18 09:45:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-20 17:00:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-20 17:00:00');

-- dnt_joana_01 (posto: São Luiz Star / residência: adr_usr_joana)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_joana01_exame', 'dnt_joana_01', 'adr_dpt_saoluiz', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-11-16 09:00:00', TIMESTAMP '2025-11-15 10:10:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-16 10:00:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-16 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_joana01_kit', 'dnt_joana_01', 'adr_usr_joana', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-11-18 10:00:00', TIMESTAMP '2025-11-16 10:05:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-18 10:30:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-18 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_joana01_coleta', 'dnt_joana_01', 'adr_usr_joana', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-11-22 09:00:00', TIMESTAMP '2025-11-18 10:35:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-22 09:40:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-22 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_joana01_analise', 'dnt_joana_01', 'adr_dpt_saoluiz', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-11-24 15:00:00', TIMESTAMP '2025-11-22 09:45:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-25 16:30:00', 'usr_nurse_juliana', TIMESTAMP '2025-11-25 16:30:00');

-- dnt_beatriz_01 (posto: Hospital Ipiranga, com coleta domiciliar / residência: adr_usr_beatriz)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_beatriz01_exame', 'dnt_beatriz_01', 'adr_dpt_ipiranga', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-12-02 09:00:00', TIMESTAMP '2025-12-01 08:40:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-02 10:00:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-02 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_beatriz01_kit', 'dnt_beatriz_01', 'adr_usr_beatriz', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-12-03 10:00:00', TIMESTAMP '2025-12-02 10:05:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-03 10:30:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-03 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_beatriz01_coleta', 'dnt_beatriz_01', 'adr_usr_beatriz', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-12-08 09:00:00', TIMESTAMP '2025-12-03 10:35:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-08 09:40:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-08 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_beatriz01_analise', 'dnt_beatriz_01', 'adr_dpt_ipiranga', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-12-10 15:00:00', TIMESTAMP '2025-12-08 09:45:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-10 18:00:00', 'usr_nurse_patricia', TIMESTAMP '2025-12-10 18:00:00');

-- dnt_paula_01 (posto: Amparo Maternal / residência: adr_usr_paula)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_paula01_exame', 'dnt_paula_01', 'adr_dpt_amparo', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-12-06 09:00:00', TIMESTAMP '2025-12-05 13:10:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-06 10:00:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-06 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_paula01_kit', 'dnt_paula_01', 'adr_usr_paula', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-12-08 10:00:00', TIMESTAMP '2025-12-06 10:05:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-08 10:30:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-08 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_paula01_coleta', 'dnt_paula_01', 'adr_usr_paula', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-12-14 09:00:00', TIMESTAMP '2025-12-08 10:35:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-14 09:40:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-14 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_paula01_analise', 'dnt_paula_01', 'adr_dpt_amparo', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-12-17 15:00:00', TIMESTAMP '2025-12-14 09:45:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-18 12:00:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-18 12:00:00');

-- dnt_gabriela_01 (posto: Anália Franco / residência: adr_usr_gabriela)
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_gabriela01_exame', 'dnt_gabriela_01', 'adr_dpt_analia', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2025-12-11 09:00:00', TIMESTAMP '2025-12-10 09:10:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-11 10:00:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-11 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_gabriela01_kit', 'dnt_gabriela_01', 'adr_usr_gabriela', 'Entregar kit de ordenha', 'Kit de ordenha entregue na residência da doadora.', 'done', TIMESTAMP '2025-12-13 10:00:00', TIMESTAMP '2025-12-11 10:05:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-13 10:30:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-13 10:30:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_gabriela01_coleta', 'dnt_gabriela_01', 'adr_usr_gabriela', 'Coletar leite', 'Coleta do leite ordenhado realizada em domicílio.', 'done', TIMESTAMP '2025-12-18 09:00:00', TIMESTAMP '2025-12-13 10:35:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-18 09:40:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-18 09:40:00');
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_gabriela01_analise', 'dnt_gabriela_01', 'adr_dpt_analia', 'Análise de leite', 'Leite aprovado nos testes de qualidade e segurança.', 'done', TIMESTAMP '2025-12-21 15:00:00', TIMESTAMP '2025-12-18 09:45:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-22 17:30:00', 'usr_nurse_fernanda', TIMESTAMP '2025-12-22 17:30:00');

-- doações em andamento: exame concluído, kit pendente
INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_maria02_exame', 'dnt_maria_02', 'adr_dpt_amparo', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-06 09:00:00', TIMESTAMP '2026-07-05 09:10:00', 'usr_nurse_juliana', TIMESTAMP '2026-07-06 10:00:00', 'usr_nurse_juliana', TIMESTAMP '2026-07-06 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('stp_maria02_kit', 'dnt_maria_02', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-06 10:05:00', 'usr_nurse_juliana');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_camila01_exame', 'dnt_camila_01', 'adr_dpt_analia', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-06-21 09:00:00', TIMESTAMP '2026-06-20 11:10:00', 'usr_nurse_patricia', TIMESTAMP '2026-06-21 10:00:00', 'usr_nurse_patricia', TIMESTAMP '2026-06-21 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('stp_camila01_kit', 'dnt_camila_01', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-06-21 10:05:00', 'usr_nurse_patricia');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_larissa01_exame', 'dnt_larissa_01', 'adr_dpt_santacasa', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-11 09:00:00', TIMESTAMP '2026-07-10 09:40:00', 'usr_nurse_patricia', TIMESTAMP '2026-07-11 10:00:00', 'usr_nurse_patricia', TIMESTAMP '2026-07-11 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('stp_larissa01_kit', 'dnt_larissa_01', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-11 10:05:00', 'usr_nurse_patricia');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_amanda01_exame', 'dnt_amanda_01', 'adr_dpt_saoluiz', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-16 09:00:00', TIMESTAMP '2026-07-15 14:10:00', 'usr_nurse_fernanda', TIMESTAMP '2026-07-16 10:00:00', 'usr_nurse_fernanda', TIMESTAMP '2026-07-16 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('stp_amanda01_kit', 'dnt_amanda_01', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-16 10:05:00', 'usr_nurse_fernanda');

INSERT INTO donation_step (id_donation_step, id_donation, id_address, name, description, status, set_date, created_at, created_by, updated_at, updated_by, completed_at)
VALUES ('stp_rafaela01_exame', 'dnt_rafaela_01', 'adr_dpt_ipiranga', 'Exame de sangue', 'Triagem sorológica realizada no posto de coleta.', 'done', TIMESTAMP '2026-07-21 09:00:00', TIMESTAMP '2026-07-20 10:10:00', 'usr_nurse_fernanda', TIMESTAMP '2026-07-21 10:00:00', 'usr_nurse_fernanda', TIMESTAMP '2026-07-21 10:00:00');
INSERT INTO donation_step (id_donation_step, id_donation, name, description, status, created_at, created_by)
VALUES ('stp_rafaela01_kit', 'dnt_rafaela_01', 'Entregar kit de ordenha', 'Kit de ordenha aguardando agendamento de entrega.', 'pending', TIMESTAMP '2026-07-21 10:05:00', 'usr_nurse_fernanda');

--------------------------------------------------------------------------------
-- DONATION STEP TIMELINE (1 registro de histórico por etapa: 30)
--------------------------------------------------------------------------------
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_maria01_exame', 'stp_maria01_exame', 'adr_dpt_amparo', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-11-11 10:00:00', TIMESTAMP '2025-11-11 10:00:00', 'usr_nurse_juliana');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_maria01_kit', 'stp_maria01_kit', 'adr_usr_maria', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-11-13 10:30:00', TIMESTAMP '2025-11-13 10:30:00', 'usr_nurse_juliana');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_maria01_coleta', 'stp_maria01_coleta', 'adr_usr_maria', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-11-18 09:40:00', TIMESTAMP '2025-11-18 09:40:00', 'usr_nurse_juliana');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_maria01_analise', 'stp_maria01_analise', 'adr_dpt_amparo', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-11-20 17:00:00', TIMESTAMP '2025-11-20 17:00:00', 'usr_nurse_juliana');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_joana01_exame', 'stp_joana01_exame', 'adr_dpt_saoluiz', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-11-16 10:00:00', TIMESTAMP '2025-11-16 10:00:00', 'usr_nurse_juliana');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_joana01_kit', 'stp_joana01_kit', 'adr_usr_joana', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-11-18 10:30:00', TIMESTAMP '2025-11-18 10:30:00', 'usr_nurse_juliana');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_joana01_coleta', 'stp_joana01_coleta', 'adr_usr_joana', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-11-22 09:40:00', TIMESTAMP '2025-11-22 09:40:00', 'usr_nurse_juliana');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_joana01_analise', 'stp_joana01_analise', 'adr_dpt_saoluiz', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-11-25 16:30:00', TIMESTAMP '2025-11-25 16:30:00', 'usr_nurse_juliana');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_beatriz01_exame', 'stp_beatriz01_exame', 'adr_dpt_ipiranga', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-12-02 10:00:00', TIMESTAMP '2025-12-02 10:00:00', 'usr_nurse_patricia');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_beatriz01_kit', 'stp_beatriz01_kit', 'adr_usr_beatriz', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-12-03 10:30:00', TIMESTAMP '2025-12-03 10:30:00', 'usr_nurse_patricia');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_beatriz01_coleta', 'stp_beatriz01_coleta', 'adr_usr_beatriz', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-12-08 09:40:00', TIMESTAMP '2025-12-08 09:40:00', 'usr_nurse_patricia');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_beatriz01_analise', 'stp_beatriz01_analise', 'adr_dpt_ipiranga', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-12-10 18:00:00', TIMESTAMP '2025-12-10 18:00:00', 'usr_nurse_patricia');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_paula01_exame', 'stp_paula01_exame', 'adr_dpt_amparo', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-12-06 10:00:00', TIMESTAMP '2025-12-06 10:00:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_paula01_kit', 'stp_paula01_kit', 'adr_usr_paula', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-12-08 10:30:00', TIMESTAMP '2025-12-08 10:30:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_paula01_coleta', 'stp_paula01_coleta', 'adr_usr_paula', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-12-14 09:40:00', TIMESTAMP '2025-12-14 09:40:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_paula01_analise', 'stp_paula01_analise', 'adr_dpt_amparo', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-12-18 12:00:00', TIMESTAMP '2025-12-18 12:00:00', 'usr_nurse_fernanda');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_gabriela01_exame', 'stp_gabriela01_exame', 'adr_dpt_analia', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2025-12-11 10:00:00', TIMESTAMP '2025-12-11 10:00:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_gabriela01_kit', 'stp_gabriela01_kit', 'adr_usr_gabriela', 'Kit de ordenha entregue e recebido pela doadora.', 'done', TIMESTAMP '2025-12-13 10:30:00', TIMESTAMP '2025-12-13 10:30:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_gabriela01_coleta', 'stp_gabriela01_coleta', 'adr_usr_gabriela', 'Leite coletado em domicílio e encaminhado ao posto.', 'done', TIMESTAMP '2025-12-18 09:40:00', TIMESTAMP '2025-12-18 09:40:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_gabriela01_analise', 'stp_gabriela01_analise', 'adr_dpt_analia', 'Leite aprovado na análise laboratorial de qualidade.', 'done', TIMESTAMP '2025-12-22 17:30:00', TIMESTAMP '2025-12-22 17:30:00', 'usr_nurse_fernanda');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_maria02_exame', 'stp_maria02_exame', 'adr_dpt_amparo', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-06 10:00:00', TIMESTAMP '2026-07-06 10:00:00', 'usr_nurse_juliana');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('tml_maria02_kit', 'stp_maria02_kit', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-06 10:05:00', 'usr_nurse_juliana');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_camila01_exame', 'stp_camila01_exame', 'adr_dpt_analia', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-06-21 10:00:00', TIMESTAMP '2026-06-21 10:00:00', 'usr_nurse_patricia');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('tml_camila01_kit', 'stp_camila01_kit', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-06-21 10:05:00', 'usr_nurse_patricia');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_larissa01_exame', 'stp_larissa01_exame', 'adr_dpt_santacasa', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-11 10:00:00', TIMESTAMP '2026-07-11 10:00:00', 'usr_nurse_patricia');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('tml_larissa01_kit', 'stp_larissa01_kit', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-11 10:05:00', 'usr_nurse_patricia');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_amanda01_exame', 'stp_amanda01_exame', 'adr_dpt_saoluiz', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-16 10:00:00', TIMESTAMP '2026-07-16 10:00:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('tml_amanda01_kit', 'stp_amanda01_kit', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-16 10:05:00', 'usr_nurse_fernanda');

INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, id_address, description, status, set_date, created_at, created_by)
VALUES ('tml_rafaela01_exame', 'stp_rafaela01_exame', 'adr_dpt_ipiranga', 'Resultado do exame de sangue liberado sem alterações.', 'done', TIMESTAMP '2026-07-21 10:00:00', TIMESTAMP '2026-07-21 10:00:00', 'usr_nurse_fernanda');
INSERT INTO donation_step_timeline (id_donation_step_timeline, id_donation_step, description, status, created_at, created_by)
VALUES ('tml_rafaela01_kit', 'stp_rafaela01_kit', 'Aguardando confirmação de horário para entrega do kit.', 'pending', TIMESTAMP '2026-07-21 10:05:00', 'usr_nurse_fernanda');

--------------------------------------------------------------------------------
-- JOB (10: visitas domiciliares de entrega de kit e coleta de leite)
--------------------------------------------------------------------------------
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_maria01_kit', 'usr_nurse_juliana', 'stp_maria01_kit', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-11-13 10:00:00', TIMESTAMP '2025-11-11 10:05:00', 'usr_nurse_juliana');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_maria01_coleta', 'usr_nurse_juliana', 'stp_maria01_coleta', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-11-18 09:00:00', TIMESTAMP '2025-11-13 10:35:00', 'usr_nurse_juliana');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_joana01_kit', 'usr_nurse_juliana', 'stp_joana01_kit', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-11-18 10:00:00', TIMESTAMP '2025-11-16 10:05:00', 'usr_nurse_juliana');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_joana01_coleta', 'usr_nurse_juliana', 'stp_joana01_coleta', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-11-22 09:00:00', TIMESTAMP '2025-11-18 10:35:00', 'usr_nurse_juliana');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_beatriz01_kit', 'usr_nurse_patricia', 'stp_beatriz01_kit', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-12-03 10:00:00', TIMESTAMP '2025-12-02 10:05:00', 'usr_nurse_patricia');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_beatriz01_coleta', 'usr_nurse_patricia', 'stp_beatriz01_coleta', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-12-08 09:00:00', TIMESTAMP '2025-12-03 10:35:00', 'usr_nurse_patricia');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_paula01_kit', 'usr_nurse_fernanda', 'stp_paula01_kit', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-12-08 10:00:00', TIMESTAMP '2025-12-06 10:05:00', 'usr_nurse_fernanda');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_paula01_coleta', 'usr_nurse_fernanda', 'stp_paula01_coleta', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-12-14 09:00:00', TIMESTAMP '2025-12-08 10:35:00', 'usr_nurse_fernanda');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_gabriela01_kit', 'usr_nurse_fernanda', 'stp_gabriela01_kit', 'done', 'Entrega de kit de ordenha', 'Visita para entrega do kit de ordenha na residência da doadora.', TIMESTAMP '2025-12-13 10:00:00', TIMESTAMP '2025-12-11 10:05:00', 'usr_nurse_fernanda');
INSERT INTO job (id_job, id_user, id_step, status, name, description, date_set, created_at, created_by)
VALUES ('job_gabriela01_coleta', 'usr_nurse_fernanda', 'stp_gabriela01_coleta', 'done', 'Coleta de leite em domicílio', 'Visita para recolhimento do leite ordenhado.', TIMESTAMP '2025-12-18 09:00:00', TIMESTAMP '2025-12-13 10:35:00', 'usr_nurse_fernanda');

--------------------------------------------------------------------------------
-- CONSENT LOG (15: um por usuário)
--------------------------------------------------------------------------------
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_adm_carla', 'usr_adm_carla', 'v1.0', TIMESTAMP '2025-01-05 09:00:00', '200.150.10.11', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_adm_bruno', 'usr_adm_bruno', 'v1.0', TIMESTAMP '2025-01-05 09:15:00', '200.150.10.12', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_nurse_juliana', 'usr_nurse_juliana', 'v1.0', TIMESTAMP '2025-01-10 10:00:00', '200.150.10.13', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_nurse_patricia', 'usr_nurse_patricia', 'v1.0', TIMESTAMP '2025-01-10 10:15:00', '200.150.10.14', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_nurse_fernanda', 'usr_nurse_fernanda', 'v1.0', TIMESTAMP '2025-01-12 11:00:00', '200.150.10.15', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) NutrizWeb/1.0');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_maria_silva', 'usr_maria_silva', 'v1.1', TIMESTAMP '2025-02-01 08:30:00', '187.10.20.1', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_joana_souza', 'usr_joana_souza', 'v1.1', TIMESTAMP '2025-02-03 09:10:00', '187.10.20.2', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_camila_lima', 'usr_camila_lima', 'v1.1', TIMESTAMP '2025-02-10 14:20:00', '187.10.20.3', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_beatriz_alves', 'usr_beatriz_alves', 'v1.1', TIMESTAMP '2025-02-15 16:00:00', '187.10.20.4', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_larissa_costa', 'usr_larissa_costa', 'v1.1', TIMESTAMP '2025-03-01 08:00:00', '187.10.20.5', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_paula_martins', 'usr_paula_martins', 'v1.1', TIMESTAMP '2025-03-05 09:45:00', '187.10.20.6', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_amanda_rocha', 'usr_amanda_rocha', 'v1.1', TIMESTAMP '2025-03-10 11:30:00', '187.10.20.7', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_gabriela_teixeira', 'usr_gabriela_teixeira', 'v1.1', TIMESTAMP '2025-03-15 13:00:00', '187.10.20.8', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_rafaela_pinto', 'usr_rafaela_pinto', 'v1.1', TIMESTAMP '2025-03-20 15:15:00', '187.10.20.9', 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_5) NutrizApp/2.3');
INSERT INTO consent_log (id_consent_log, id_user, terms_version, accepted_at, ip_address, user_agent)
VALUES ('cst_usr_debora_nunes', 'usr_debora_nunes', 'v1.1', TIMESTAMP '2025-04-02 10:00:00', '187.10.20.10', 'Mozilla/5.0 (Linux; Android 14) NutrizApp/2.3');

--------------------------------------------------------------------------------
-- USER BABY (8: filhos das doadoras, motivo da amamentação/doação)
--------------------------------------------------------------------------------
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_maria_sofia', 'usr_maria_silva', 'Sofia', TIMESTAMP '2025-01-10 00:00:00', TIMESTAMP '2025-02-01 08:32:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_maria_miguel', 'usr_maria_silva', 'Miguel', TIMESTAMP '2025-01-10 00:00:00', TIMESTAMP '2025-02-01 08:33:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_joana_helena', 'usr_joana_souza', 'Helena', TIMESTAMP '2025-03-22 00:00:00', TIMESTAMP '2025-02-03 09:12:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_camila_davi', 'usr_camila_lima', 'Davi', TIMESTAMP '2025-05-02 00:00:00', TIMESTAMP '2025-02-10 14:22:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_beatriz_laura', 'usr_beatriz_alves', 'Laura', TIMESTAMP '2024-11-15 00:00:00', TIMESTAMP '2025-02-15 16:02:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_larissa_pedro', 'usr_larissa_costa', 'Pedro', TIMESTAMP '2025-02-08 00:00:00', TIMESTAMP '2025-03-01 08:02:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_paula_isabela', 'usr_paula_martins', 'Isabela', TIMESTAMP '2025-06-19 00:00:00', TIMESTAMP '2025-03-05 09:47:00');
INSERT INTO user_baby (id_user_baby, id_user, name, birth_date, created_at)
VALUES ('bby_gabriela_enzo', 'usr_gabriela_teixeira', 'Enzo', TIMESTAMP '2025-04-30 00:00:00', TIMESTAMP '2025-03-15 13:02:00');

--------------------------------------------------------------------------------
-- KB CHUNKS (6: base de conhecimento do assistente de IA)
-- Os embeddings de 384 dimensões são gerados em PL/SQL (impraticável digitar
-- 384 números manualmente); o conteúdo textual é real e coerente com o domínio.
--------------------------------------------------------------------------------
DECLARE
  TYPE t_chunk IS RECORD (
    id        VARCHAR2(36),
    source    VARCHAR2(100),
    content   VARCHAR2(500),
    metadata  VARCHAR2(200)
  );
  TYPE t_chunks IS TABLE OF t_chunk;
  v_chunks t_chunks := t_chunks(
    t_chunk('kb_faq_elegibilidade', 'faq_doacao_leite.pdf', 'Podem doar leite humano mulheres saudaveis, nao fumantes, que nao usam determinados medicamentos e que produzem leite excedente apos amamentar o proprio bebe.', '{"topic":"elegibilidade","lang":"pt-BR"}'),
    t_chunk('kb_faq_higiene', 'manual_banco_leite.pdf', 'Antes da ordenha, a doadora deve lavar as maos e os seios com agua e sabao neutro, e utilizar frascos esterilizados fornecidos pelo banco de leite.', '{"topic":"higiene","lang":"pt-BR"}'),
    t_chunk('kb_faq_armazenamento', 'manual_banco_leite.pdf', 'O leite ordenhado deve ser armazenado em freezer, identificado com data e hora da coleta, podendo ser mantido por ate 15 dias antes da coleta pelo banco de leite.', '{"topic":"armazenamento","lang":"pt-BR"}'),
    t_chunk('kb_faq_exames', 'faq_doacao_leite.pdf', 'A doadora passa por uma triagem com exame de sangue para descartar doencas infectocontagiosas antes de iniciar as doacoes regulares.', '{"topic":"exames","lang":"pt-BR"}'),
    t_chunk('kb_faq_beneficios', 'guia_amamentacao.pdf', 'O leite humano doado e destinado prioritariamente a bebes prematuros e internados em UTI neonatal, auxiliando no fortalecimento do sistema imunologico.', '{"topic":"beneficios","lang":"pt-BR"}'),
    t_chunk('kb_faq_kit_coleta', 'guia_amamentacao.pdf', 'O kit de ordenha, com frascos e etiquetas, e entregue gratuitamente na casa da doadora ou em um dos postos de coleta cadastrados no Nutriz.', '{"topic":"kit_coleta","lang":"pt-BR"}')
  );
  v_vec VARCHAR2(32767);
BEGIN
  FOR i IN 1 .. v_chunks.COUNT LOOP
    v_vec := '[';
    FOR j IN 1 .. 384 LOOP
      v_vec := v_vec || TO_CHAR(ROUND(DBMS_RANDOM.VALUE(-1, 1), 4));
      IF j < 384 THEN
        v_vec := v_vec || ',';
      END IF;
    END LOOP;
    v_vec := v_vec || ']';

    INSERT INTO kb_chunks (id, source, content, embedding, metadata)
    VALUES (
      v_chunks(i).id,
      v_chunks(i).source,
      v_chunks(i).content,
      TO_VECTOR(v_vec, 384, FLOAT32),
      v_chunks(i).metadata
    );
  END LOOP;
END;
/

--------------------------------------------------------------------------------
-- CONVERSATIONS (8: doadoras conversando com o assistente de IA)
--------------------------------------------------------------------------------
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_maria_01', 'usr_maria_silva', TIMESTAMP '2026-07-05 08:00:00', TIMESTAMP '2026-07-05 08:06:00', 'Dúvidas sobre nova doação e reagendamento de etapas.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_joana_01', 'usr_joana_souza', TIMESTAMP '2026-07-06 09:00:00', TIMESTAMP '2026-07-06 09:03:00', 'Dúvida sobre prazo de armazenamento do leite no freezer.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_camila_01', 'usr_camila_lima', TIMESTAMP '2026-06-19 10:00:00', TIMESTAMP '2026-06-19 10:08:00', 'Primeira doadora perguntando como iniciar o processo de doação.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_beatriz_01', 'usr_beatriz_alves', TIMESTAMP '2025-12-11 11:00:00', TIMESTAMP '2025-12-11 11:02:00', 'Consulta sobre resultado da análise de leite.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_larissa_01', 'usr_larissa_costa', TIMESTAMP '2026-07-09 15:00:00', TIMESTAMP '2026-07-09 15:02:00', 'Solicitação de alteração de endereço para coleta.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_amanda_01', 'usr_amanda_rocha', TIMESTAMP '2026-07-16 16:00:00', TIMESTAMP '2026-07-16 16:02:00', 'Dúvida sobre prazo de entrega do kit após exame.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_rafaela_01', 'usr_rafaela_pinto', TIMESTAMP '2026-07-20 17:00:00', TIMESTAMP '2026-07-20 17:02:00', 'Pergunta sobre coleta domiciliar no Hospital Ipiranga.');
INSERT INTO conversations (id, user_id, started_at, last_message_at, summary)
VALUES ('conv_debora_01', 'usr_debora_nunes', TIMESTAMP '2026-08-01 12:00:00', TIMESTAMP '2026-08-01 12:03:00', 'Usuária ainda não doadora perguntando sobre requisitos de elegibilidade.');

--------------------------------------------------------------------------------
-- MESSAGES (20: alternando usuário/assistente)
--------------------------------------------------------------------------------
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_maria01_1', 'conv_maria_01', 'user', 'Já fiz uma doação em novembro, posso doar de novo?', 14, TIMESTAMP '2026-07-05 08:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_maria01_2', 'conv_maria_01', 'assistant', 'Sim! Você pode doar novamente após passar por uma nova avaliação de saúde. Vou agendar o exame de sangue para você.', 28, TIMESTAMP '2026-07-05 08:02:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_maria01_3', 'conv_maria_01', 'user', 'Perfeito, quando será a entrega do kit?', 10, TIMESTAMP '2026-07-05 08:05:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_maria01_4', 'conv_maria_01', 'assistant', 'A entrega do kit de ordenha está aguardando agendamento; uma de nossas enfermeiras entrará em contato para confirmar o horário.', 26, TIMESTAMP '2026-07-05 08:06:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_joana01_1', 'conv_joana_01', 'user', 'Quanto tempo posso guardar o leite ordenhado no freezer antes da coleta?', 16, TIMESTAMP '2026-07-06 09:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_joana01_2', 'conv_joana_01', 'assistant', 'O leite pode ser armazenado no freezer por até 15 dias antes da coleta pelo banco de leite, sempre identificado com data e hora.', 24, TIMESTAMP '2026-07-06 09:03:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_camila01_1', 'conv_camila_01', 'user', 'Estou amamentando meu filho Davi e tenho leite excedente, como funciona a doação?', 19, TIMESTAMP '2026-06-19 10:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_camila01_2', 'conv_camila_01', 'assistant', 'Que ótimo, Camila! Após um exame de sangue de triagem, você recebe um kit de ordenha e pode começar a doar. Deseja que eu inicie o processo?', 30, TIMESTAMP '2026-06-19 10:04:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_camila01_3', 'conv_camila_01', 'user', 'Sim, pode iniciar.', 6, TIMESTAMP '2026-06-19 10:06:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_camila01_4', 'conv_camila_01', 'assistant', 'Prontinho! Já registrei sua doação e agendei o exame de sangue inicial.', 15, TIMESTAMP '2026-06-19 10:08:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_beatriz01_1', 'conv_beatriz_01', 'user', 'Recebi a notificação de que minha doação foi analisada, qual foi o resultado?', 17, TIMESTAMP '2025-12-11 11:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_beatriz01_2', 'conv_beatriz_01', 'assistant', 'Sua doação de 1200ml foi aprovada na análise de leite e já está disponível para os bebês atendidos pelo Hospital Ipiranga.', 27, TIMESTAMP '2025-12-11 11:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_larissa01_1', 'conv_larissa_01', 'user', 'Preciso trocar o endereço para a coleta do kit, como faço?', 14, TIMESTAMP '2026-07-09 15:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_larissa01_2', 'conv_larissa_01', 'assistant', 'Claro, você pode atualizar seu endereço no aplicativo e a próxima etapa será ajustada automaticamente.', 20, TIMESTAMP '2026-07-09 15:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_amanda01_1', 'conv_amanda_01', 'user', 'Fiz o exame de sangue, quando recebo o kit de ordenha?', 13, TIMESTAMP '2026-07-16 16:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_amanda01_2', 'conv_amanda_01', 'assistant', 'Seu exame foi processado e o kit de ordenha será entregue em até 3 dias úteis no endereço cadastrado.', 21, TIMESTAMP '2026-07-16 16:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_rafaela01_1', 'conv_rafaela_01', 'user', 'O Hospital Ipiranga faz coleta em domicílio?', 10, TIMESTAMP '2026-07-20 17:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_rafaela01_2', 'conv_rafaela_01', 'assistant', 'Sim! O Banco de Leite Humano do Hospital Ipiranga oferece coleta domiciliar, ideal para o seu caso.', 19, TIMESTAMP '2026-07-20 17:02:00');

INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_debora01_1', 'conv_debora_01', 'user', 'Ainda não doei, quais são os requisitos para me tornar doadora?', 15, TIMESTAMP '2026-08-01 12:00:00');
INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at)
VALUES ('msg_debora01_2', 'conv_debora_01', 'assistant', 'Você precisa ser uma mulher saudável, não fumante, sem uso de determinados medicamentos, e produzir leite excedente após amamentar seu bebê. Deseja iniciar a triagem?', 32, TIMESTAMP '2026-08-01 12:03:00');

--------------------------------------------------------------------------------
-- LLM AUDIT (10: uma auditoria por resposta do assistente)
--------------------------------------------------------------------------------
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_maria01_2', 'usr_maria_silva', 'conv_maria_01', 'msg_maria01_2', '{"system":"Assistente Nutriz","user":"Já fiz uma doação em novembro, posso doar de novo?"}', '["kb_faq_elegibilidade"]', 'anthropic', 'claude-sonnet-5', 120, 28, 850, FALSE, 'sess_maria_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_maria01_4', 'usr_maria_silva', 'conv_maria_01', 'msg_maria01_4', '{"system":"Assistente Nutriz","user":"Perfeito, quando será a entrega do kit?"}', '["kb_faq_kit_coleta"]', 'anthropic', 'claude-sonnet-5', 95, 26, 780, FALSE, 'sess_maria_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_joana01_2', 'usr_joana_souza', 'conv_joana_01', 'msg_joana01_2', '{"system":"Assistente Nutriz","user":"Quanto tempo posso guardar o leite ordenhado no freezer antes da coleta?"}', '["kb_faq_armazenamento"]', 'anthropic', 'claude-sonnet-5', 110, 24, 720, FALSE, 'sess_joana_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_camila01_2', 'usr_camila_lima', 'conv_camila_01', 'msg_camila01_2', '{"system":"Assistente Nutriz","user":"Estou amamentando meu filho Davi e tenho leite excedente, como funciona a doação?"}', '["kb_faq_elegibilidade","kb_faq_kit_coleta"]', 'anthropic', 'claude-sonnet-5', 140, 30, 910, FALSE, 'sess_camila_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_camila01_4', 'usr_camila_lima', 'conv_camila_01', 'msg_camila01_4', '{"system":"Assistente Nutriz","user":"Sim, pode iniciar."}', NULL, 'anthropic', 'claude-sonnet-5', 60, 15, 500, FALSE, 'sess_camila_01', 'start_donation');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_beatriz01_2', 'usr_beatriz_alves', 'conv_beatriz_01', 'msg_beatriz01_2', '{"system":"Assistente Nutriz","user":"Recebi a notificação de que minha doação foi analisada, qual foi o resultado?"}', NULL, 'anthropic', 'claude-sonnet-5', 105, 27, 760, FALSE, 'sess_beatriz_01', 'answer_status');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_larissa01_2', 'usr_larissa_costa', 'conv_larissa_01', 'msg_larissa01_2', '{"system":"Assistente Nutriz","user":"Preciso trocar o endereço para a coleta do kit, como faço?"}', NULL, 'anthropic', 'claude-sonnet-5', 90, 20, 640, FALSE, 'sess_larissa_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_amanda01_2', 'usr_amanda_rocha', 'conv_amanda_01', 'msg_amanda01_2', '{"system":"Assistente Nutriz","user":"Fiz o exame de sangue, quando recebo o kit de ordenha?"}', '["kb_faq_kit_coleta"]', 'anthropic', 'claude-sonnet-5', 100, 21, 700, FALSE, 'sess_amanda_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_rafaela01_2', 'usr_rafaela_pinto', 'conv_rafaela_01', 'msg_rafaela01_2', '{"system":"Assistente Nutriz","user":"O Hospital Ipiranga faz coleta em domicílio?"}', NULL, 'anthropic', 'claude-sonnet-5', 85, 19, 610, FALSE, 'sess_rafaela_01', 'answer_faq');
INSERT INTO llm_audit (id, user_id, conversation_id, message_id, prompt_full, chunks_used, llm_provider, llm_model, tokens_input, tokens_output, latency_ms, is_anonymous, session_id, action_emitted)
VALUES ('aud_debora01_2', 'usr_debora_nunes', 'conv_debora_01', 'msg_debora01_2', '{"system":"Assistente Nutriz","user":"Ainda não doei, quais são os requisitos para me tornar doadora?"}', '["kb_faq_elegibilidade"]', 'anthropic', 'claude-sonnet-5', 130, 32, 880, FALSE, 'sess_debora_01', 'answer_faq');

--------------------------------------------------------------------------------
-- ALEMBIC VERSION (controle de migrations do serviço de IA)
--------------------------------------------------------------------------------
INSERT INTO alembic_version (version_num) VALUES ('d4e91a7c22b0');

COMMIT;
