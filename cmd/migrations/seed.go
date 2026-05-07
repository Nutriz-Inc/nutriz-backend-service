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
    'senha_hash_1',
    2.5,
    NOW(),
    'usr_2veL1FPpuXxUaZcFaEC57BfpcKE',
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
);`
