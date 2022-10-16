-- +goose Up
-- +goose StatementBegin

create table if not exists exchange_rates
(
    id            integer generated always as identity
        primary key,
    currency_code text           not null,
    rate          numeric(10, 7) not null,
    created_at    timestamp      not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE exchange_rates;

-- +goose StatementEnd
