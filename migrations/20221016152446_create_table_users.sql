-- +goose Up
-- +goose StatementBegin

create table if not exists users
(
    user_id       integer   not null
        primary key,
    currency_code text,
    monthly_limit integer,
    updated_at    timestamp not null
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE users;

-- +goose StatementEnd
