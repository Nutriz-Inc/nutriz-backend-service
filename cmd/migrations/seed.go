package main

var seed = `
INSERT INTO "user" (
    id_user,
    internal_identifier,
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
    '234567898765435',
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
    '234567898765436',
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
    quantity_donated,
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
    'adr_01JTG8J5F6W9K2M4P7Q1X8Y3ZAL',
    NULL,
    'dpt_2veL1FPpuXxUaZcFaEC57BfpcKE',
    '01310100',
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
    'adr_01JTG8K8N4P2R6T9V1X3Y5Z7BCD',
    NULL,
    'dpt_2veL1FPpuXxUaZcFaEC57BfpcMA',
    '22250040',
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
),
(
    'adr_01JTG8K8N4P2R6T9V1X3Y5Z7MIP',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    NULL,
    '09415987',
    'Rua Ribeiros dos Mares',
    '75',
    'São Bernardo do Campo',
    'SP',
    'Centro',
    'Próximo ao Hospital São Paulo',
    -27.923916,
    -46.172981,
    NOW(),
    NULL,
    NULL,
    NULL,
    NULL
),
(
    'adr_01JTX0H1V8N5Q3W7E2R4T6Y8UPD',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    NULL,
    '01545000',
    'Rua Teste Atualizacao',
    '101',
    'São Paulo',
    'SP',
    'Vila Mariana',
    'Casa 1',
    -23.589523,
    -46.637636,
    NOW(),
    NULL,
    NULL,
    NULL,
    NULL
),
(
    'adr_01JTX0H1V8N5Q3W7E2R4T6Y8MBL',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKL',
    NULL,
    '01545001',
    'Rua Teste Atualizacao 2',
    '102',
    'São Paulo',
    'SP',
    'Vila Madalena',
    'Casa 5',
    -22.589923,
    -41.637736,
    NOW(),
    NULL,
    NULL,
    NULL,
    NULL
);

INSERT INTO user_baby (
    id_user_baby,
    id_user,
    name,
    birth_date,
    created_at,
    updated_at,
    removed_at
)
VALUES
(
    'usb_01JTG8J5F6W9K2M4P7Q1X8Y3ZAL',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    'Miguel',
    '2025-01-15 08:30:00',
    NOW(),
    NULL,
    NULL
),
(
    'usb_01JTG8K8N4P2R6T9V1X3Y5Z7BCD',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    'Helena',
    '2024-11-02 14:10:00',
    NOW(),
    NULL,
    NULL
),
(
    'usb_01JTG8K8N4P2R6T9V1X3Y5Z7MIP',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKL',
    'Arthur',
    '2025-03-20 06:45:00',
    NOW(),
    NULL,
    NULL
);

INSERT INTO donation_step (
    id_step_donation,
    id_donation,
    name,
    description,
    status,
    set_date,
    created_at,
    updated_at,
    completed_at
) VALUES
(
    'dst_2veL1FPpuXxUaZcFaEC57BfpcAA',
    'don_2veL1FPpuXxUaZcFaEC57BfpcKE',
    'Coleta inicial',
    'Coleta agendada com o doador',
    'pending',
    NOW() + INTERVAL '2 days',
    NOW(),
    NULL,
    NULL
),
(
    'dst_2veL1FPpuXxUaZcFaEC57BfpcBB',
    'don_2veL1FPpuXxUaZcFaEC57BfpcKF',
    'Triagem',
    'Triagem do leite doado',
    'review',
    NOW() + INTERVAL '3 days',
    NOW(),
    NULL,
    NULL
);

INSERT INTO job (
    id_job,
    id_user,
    id_step,
    name,
    description,
    date_set,
    user_feedback,
    created_at,
    created_by,
    updated_at,
    updated_by,
    removed_at,
    removed_by
) VALUES
(
    'job_2veL1FPpuXxUaZcFaEC57Bfpc01',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    'dst_2veL1FPpuXxUaZcFaEC57BfpcAA',
    'Entrega ao ponto',
    'Transporte do material para o ponto de coleta',
    CURRENT_DATE + INTERVAL '2 days',
    NULL,
    NOW(),
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
    NULL,
    NULL,
    NULL,
    NULL
),
(
    'job_2veL1FPpuXxUaZcFaEC57Bfpc02',
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKL',
    'dst_2veL1FPpuXxUaZcFaEC57BfpcBB',
    'Análise laboratorial',
    'Analisar e aprovar amostras',
    CURRENT_DATE + INTERVAL '4 days',
    NULL,
    NOW(),
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKL',
    NULL,
    NULL,
    NULL,
    NULL
);
`
