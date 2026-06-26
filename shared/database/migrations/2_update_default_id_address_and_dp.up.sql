create or replace function generate_ksuid() returns text as $$
declare
    v_time timestamp with time zone := null;
    v_seconds numeric(50) := null;
    v_numeric numeric(50) := null;
    v_epoch numeric(50) := 1400000000; -- 2014-05-13T16:53:20Z
    v_base62 text := '';
    v_alphabet text[] := array[
        '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
        'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
        'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
        'U', 'V', 'W', 'X', 'Y', 'Z',
        'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
        'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
        'u', 'v', 'w', 'x', 'y', 'z'];
begin
    v_time := clock_timestamp();
    v_seconds := floor(EXTRACT(EPOCH FROM v_time)) - v_epoch;
    v_numeric := v_seconds * pow(2::numeric(50), 128)
        + ((random()::numeric(70,20) * pow(2::numeric(70,20), 48))::numeric(50) * pow(2::numeric(50), 80)::numeric(50))
        + ((random()::numeric(70,20) * pow(2::numeric(70,20), 40))::numeric(50) * pow(2::numeric(50), 40)::numeric(50))
        +  (random()::numeric(70,20) * pow(2::numeric(70,20), 40))::numeric(50);

    while v_numeric <> 0 loop
        v_base62 := v_base62 || v_alphabet[mod(v_numeric, 62)::int + 1];
        v_numeric := div(v_numeric, 62);
    end loop;
    v_base62 := reverse(v_base62);
    v_base62 := lpad(v_base62, 27, '0');

    return v_base62;
end $$ language plpgsql;

ALTER TABLE donation_point ALTER COLUMN id_donation_point SET DEFAULT 'dpt_' || generate_ksuid();
ALTER TABLE address ALTER COLUMN id_address SET DEFAULT 'adr_' || generate_ksuid();