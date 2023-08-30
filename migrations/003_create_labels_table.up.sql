CREATE TABLE labels
(
    name         VARCHAR(64) ,
    value         VARCHAR(64) ,
    pin_id    VARCHAR(64) REFERENCES pins (id),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

create index labels_pin_id_idx on labels(pin_id);
create index labels_name_idx on labels(name);
create index labels_value_idx on labels(value);