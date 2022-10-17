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

CREATE INDEX exchange_rates_currency_code_idx ON exchange_rates(currency_code);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

drop index exchange_rates_currency_code_idx;
DROP TABLE exchange_rates;

-- +goose StatementEnd
