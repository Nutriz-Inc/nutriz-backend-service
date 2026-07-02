DELETE FROM address
WHERE id_donation_point IN (
    SELECT id_donation_point FROM donation_point
    WHERE name IN (
        'Posto de Coleta de Leite Humano Amparo Maternal',
        'Banco de Leite Humano Maternidade São Luiz Star',
        'Banco de Leite Humano Rede Dor São Luiz - Unidade Anália Franco',
        'Banco de Leite Humano do Hospital Ipiranga',
        'Banco de Leite Humano da Santa Casa de São Paulo'
    )
);

DELETE FROM donation_point
WHERE name IN (
    'Posto de Coleta de Leite Humano Amparo Maternal',
    'Banco de Leite Humano Maternidade São Luiz Star',
    'Banco de Leite Humano Rede Dor São Luiz - Unidade Anália Franco',
    'Banco de Leite Humano do Hospital Ipiranga',
    'Banco de Leite Humano da Santa Casa de São Paulo'
);