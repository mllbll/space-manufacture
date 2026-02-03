-- +goose Up
create table orders (
  order_uuid text not null primary key,
  user_uuid text not null,
  part_uuids text[] not null,
  total_price real not null,
  transaction_uuid text,
  payment_method integer,
  status text not null
);

-- +goose Down
drop table orders;
