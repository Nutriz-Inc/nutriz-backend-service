INSERT INTO donation_point (name, description, has_home, phone_number, email, opening_hours)
VALUES (
    'Posto de Coleta de Leite Humano Amparo Maternal',
    'Posto de coleta de leite humano localizado no Amparo Maternal, Vila Clementino.',
    false,
    '1150898277',
    'ouvidoria@ham.spdmpais.org.br',
    'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 17:00'
);
INSERT INTO address (id_donation_point, zipcode, street, number, city, state, neighborhood, complement, latitude, longitude, created_at)
SELECT id_donation_point, '04040033', 'Rua Loefgren', '101', 'São Paulo', 'SP', 'Vila Clementino', null, -23.600032, -46.643398, now()
FROM donation_point WHERE name = 'Posto de Coleta de Leite Humano Amparo Maternal';

INSERT INTO donation_point (name, description, has_home, phone_number, email, opening_hours)
VALUES (
    'Banco de Leite Humano Maternidade São Luiz Star',
    'Banco de leite humano localizado na Maternidade São Luiz Star, Vila Olímpia.',
    false,
    '1121211349',
    'consultoriamaternidade@maternidadesaoluiz.com.br',
    'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00'
);
INSERT INTO address (id_donation_point, zipcode, street, number, city, state, neighborhood, complement, latitude, longitude, created_at)
SELECT id_donation_point, '04552050', 'Rua Helena', '29', 'São Paulo', 'SP', 'Vila Olímpia', null, -23.590756, -46.673653, now()
FROM donation_point WHERE name = 'Banco de Leite Humano Maternidade São Luiz Star';

INSERT INTO donation_point (name, description, has_home, phone_number, email, opening_hours)
VALUES (
    'Banco de Leite Humano Rede Dor São Luiz - Unidade Anália Franco',
    'Banco de leite humano localizado no Hospital e Maternidade São Luiz, Unidade Anália Franco, Tatuapé.',
    false,
    '1133861315',
    'faleconosco.sadt@saoluiz.com.br',
    'Dom-Sáb, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00'
);
INSERT INTO address (id_donation_point, zipcode, street, number, city, state, neighborhood, complement, latitude, longitude, created_at)
SELECT id_donation_point, '03313000', 'Rua Francisco Marengo', '1312', 'São Paulo', 'SP', 'Tatuapé', null, -23.548886, -46.558248, now()
FROM donation_point WHERE name = 'Banco de Leite Humano Rede Dor São Luiz - Unidade Anália Franco';

INSERT INTO donation_point (name, description, has_home, phone_number, email, opening_hours)
VALUES (
    'Banco de Leite Humano do Hospital Ipiranga',
    'Banco de leite humano localizado no Hospital Ipiranga, 8º andar.',
    true,
    '1120677866',
    'hi.cemed@gmail.com',
    null
);
INSERT INTO address (id_donation_point, zipcode, street, number, city, state, neighborhood, complement, latitude, longitude, created_at)
SELECT id_donation_point, '04262000', 'Avenida Nazaré', '28', 'São Paulo', 'SP', 'Ipiranga', '8º andar', -23.584210, -46.611595, now()
FROM donation_point WHERE name = 'Banco de Leite Humano do Hospital Ipiranga';

INSERT INTO donation_point (name, description, has_home, phone_number, email, opening_hours)
VALUES (
    'Banco de Leite Humano da Santa Casa de São Paulo',
    'Banco de leite humano localizado no Hospital Central da Santa Casa de São Paulo, Vila Buarque.',
    false,
    '1121767390',
    'ouvidoria@santacasasp.org.br',
    'Seg-Sex, Manhã: 07:00 às 11:55 | Tarde: 12:00 às 19:00'
);
INSERT INTO address (id_donation_point, zipcode, street, number, city, state, neighborhood, complement, latitude, longitude, created_at)
SELECT id_donation_point, '01221020', 'Rua Dr. Cesário Mota Júnior', '112', 'São Paulo', 'SP', 'Vila Buarque', null, -23.542626, -46.650065, now()
FROM donation_point WHERE name = 'Banco de Leite Humano da Santa Casa de São Paulo';