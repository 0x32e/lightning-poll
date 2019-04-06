create table polls(
  id bigint not null,
  status tinyint not null,
  created_at datetime not null,
  question text not null,
  expiry_seconds bigint not null,
  repay_scheme tinyint not null,
  vote_sats bigint not null,
  payout_invoice text,
  user_id bigint,

  primary key(id)
);

create table poll_options(
  id bigint not null,
  poll_id bigint not null,
  value text not null,

  primary key(id)
);

create table votes(
  id bigint not null,
  created_at datetime not null,
  poll_id bigint not null,
  option_id bigint not null,
  pay_req  text not null,
  status tinyint not null,

  primary key(id)
);
