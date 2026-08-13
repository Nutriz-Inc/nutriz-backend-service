ALTER TABLE "donation_step"
ADD COLUMN id_address VARCHAR(36),
ADD CONSTRAINT fk_donation_step_id_address
    FOREIGN KEY (id_address)
    REFERENCES "address"(id_address);

ALTER TABLE "donation_step_timeline"
ADD COLUMN id_address VARCHAR(36),
ADD CONSTRAINT fk_donation_step_timeline_id_address
    FOREIGN KEY (id_address)
    REFERENCES "address"(id_address);