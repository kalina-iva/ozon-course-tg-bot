-- +goose Up
-- +goose StatementBegin

create table if not exists expenses
(
    id         integer generated always as identity
        primary key,
    user_id    bigint    not null,
    category   text      not null,
    amount     integer   not null,
    created_at timestamp not null
);

CREATE INDEX expenses_user_id_category_idx ON expenses(user_id, category);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX expenses_user_id_category_idx;
DROP TABLE expenses;

-- +goose StatementEnd
