-- +goose Up
-- +goose StatementBegin

create table if not exists users
(
    user_id       bigint   not null
        primary key,
    currency_code text,
    monthly_limit bigint,
    updated_at    timestamp not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users;

-- +goose StatementEnd
