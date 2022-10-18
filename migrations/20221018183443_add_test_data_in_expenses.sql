-- +goose Up
-- +goose StatementBegin

INSERT INTO expenses (user_id, category, amount, created_at)
VALUES (5374344485, 'cat1', 123099, '2022-10-18 20:22:12.000000'),
       (5374344485, 'cat1', 124000, '2022-10-18 20:23:28.000000'),
       (5374344485, 'cat1', 125098, '2022-10-01 20:23:28.000000'),
       (5374344485, 'cat1', 126099, '2022-01-18 18:39:28.000000'),
       (5374344485, 'cat1', 127000, '2022-05-01 18:39:28.000000'),
       (5374344485, 'cat2', 30109, '2022-10-18 20:22:12.000000'),
       (5374344485, 'cat2', 2101000, '2022-10-18 20:23:28.000000'),
       (5374344485, 'cat2', 1790000, '2022-10-01 20:23:28.000000'),
       (5374344485, 'cat2', 709090, '2022-03-20 18:39:28.000000'),
       (5374344485, 'cat2', 1000, '2022-01-18 18:39:28.000000');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

delete from expenses where user_id = 5374344485 and category = 'cat1' and amount = 123099 and created_at = '2022-10-18 20:22:12.000000';
delete from expenses where user_id = 5374344485 and category = 'cat1' and amount = 124000 and created_at = '2022-10-18 20:23:28.000000';
delete from expenses where user_id = 5374344485 and category = 'cat1' and amount = 125098 and created_at = '2022-10-01 20:23:28.000000';
delete from expenses where user_id = 5374344485 and category = 'cat1' and amount = 126099 and created_at = '2022-01-18 18:39:28.000000';
delete from expenses where user_id = 5374344485 and category = 'cat1' and amount = 127000 and created_at = '2022-05-01 18:39:28.000000';
delete from expenses where user_id = 5374344485 and category = 'cat2' and amount = 30109 and created_at = '2022-10-18 20:22:12.000000';
delete from expenses where user_id = 5374344485 and category = 'cat2' and amount = 2101000 and created_at = '2022-10-18 20:23:28.000000';
delete from expenses where user_id = 5374344485 and category = 'cat2' and amount = 1790000 and created_at = '2022-10-01 20:23:28.000000';
delete from expenses where user_id = 5374344485 and category = 'cat2' and amount = 709090 and created_at = '2022-03-20 18:39:28.000000';
delete from expenses where user_id = 5374344485 and category = 'cat2' and amount = 1000 and created_at = '2022-01-18 18:39:28.000000';

-- +goose StatementEnd
