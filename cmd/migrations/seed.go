package main

var seed = `
INSERT INTO "user" (
    id_user,
    type,
    name,
    cpf,
    birth_date,
    phone_number,
    email,
    password,
    milk_donated,
    created_at,
    created_by,
    updated_at,
    updated_by,
    removed_at,
    removed_by
) VALUES
(
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    'common',
    'Maria Silva',
    '47046117012',
    CURRENT_DATE - INTERVAL '20 years',
    '11999999999',
    'maria@email.com',
    '2b9643d9671363af30ebed5463130ee7',
    2.5,
    NOW(),
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    NULL,
    NULL,
    NULL,
    NULL
),
(
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKL',
    'common',
    'Marta Silveira',
    '72555730028',
    CURRENT_DATE - INTERVAL '30 years',
    '11999999998',
    'marta@email.com',
    '2b9643d9671363af30ebed5463130ee7',
    5.5,
    NOW(),
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKL',
    NULL,
    NULL,
    NULL,
    NULL
);

INSERT INTO donation_point (
    id_donation_point,
    name,
    description,
    has_home,
    phone_number,
    email,
    opening_hours
) VALUES
(
    'dpt_2veL1FPpuXxUaZcFaEC57BfpcKE',
    'Ponto Solidário Centro',
    'Arrecadação de roupas e alimentos',
    true,
    '11999999999',
    'contato@pontosolidario.org',
    'Seg-Sex 08:00-18:00'
),
(
    'dpt_2veL1FPpuXxUaZcFaEC57BfpcMA',
    'Doa Fácil Zona Sul',
    NULL,
    false,
    NULL,
    'doafacil@email.com',
    'Seg-Sáb 09:00-17:00'
);

INSERT INTO donation (
    id_donation,
    is_active,
    quantity,
    user_feedback,
    created_at,
    created_by,
    updated_at,
    updated_by,
    removed_at,
    removed_by
) VALUES
(
    'don_2veL1FPpuXxUaZcFaEC57BfpcKE',
    true,
    1.5,
    'Doação realizada com sucesso',
    NOW(),
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    NULL,
    NULL,
    NULL,
    NULL
),
(
    'don_2veL1FPpuXxUaZcFaEC57BfpcKF',
    true,
    2.0,
    'Equipe muito atenciosa',
    NOW(),
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKL',
    NULL,
    NULL,
    NULL,
    NULL
);

INSERT INTO address (
    id_address,
    id_user,
    id_donation_point,
    zipcode,
    street,
    number,
    city,
    state,
    neighborhood,
    complement,
    latitude,
    longitude,
    created_at,
    updated_at,
    updated_by,
    removed_at,
    removed_by
) VALUES
(
    'adr_01JTG8J5F6W9K2M4P7Q1X8Y3ZA',
    NULL,
    'dpt_2veL1FPpuXxUaZcFaEC57BfpcKE',
    '01310-100',
    'Avenida Paulista',
    '1578',
    'São Paulo',
    'SP',
    'Bela Vista',
    'Próximo ao metrô Trianon',
    -23.561399,
    -46.656571,
    NOW(),
    NULL,
    NULL,
    NULL,
    NULL
),
(
    'adr_01JTG8K8N4P2R6T9V1X3Y5Z7BC',
    NULL,
    'dpt_2veL1FPpuXxUaZcFaEC57BfpcMA',
    '22250-040',
    'Rua Voluntários da Pátria',
    '45',
    'Rio de Janeiro',
    'RJ',
    'Copacabana',
    'Fundos',
    -22.951916,
    -43.182482,
    NOW(),
    NULL,
    NULL,
    NULL,
    NULL
);`
