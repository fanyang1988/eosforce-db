create table if not exists public.accounts
(
	name varchar(32) not null
		constraint accounts_pkey
			primary key,
	create_at timestamp,
	update_at timestamp,
	creater varchar(32),
	ref_block_num bigint,
	data jsonb
);

create table if not exists public.blocks
(
	block_num bigint not null
		constraint blocks_pkey
			primary key,
	block_id varchar(64),
	prev_block_id varchar(64),
	producer varchar(32),
	produce_time timestamp,
	transaction_merkle_root varchar(64),
	action_merkle_root varchar(64),
	new_producers varchar(64),
	num_transactions integer,
	confirmed integer
);

create index if not exists blocks_produce_time_index
	on public.blocks (produce_time);

create unique index if not exists blocks_block_id_uindex
	on public.blocks (block_id);

create table if not exists public.transfers
(
	id bigserial not null
		constraint transfer_pkey
			primary key,
	from_account varchar(32),
	to_account varchar(32),
	quantity bigint,
	token varchar(16),
	memo varchar(256),
	ref_block_num bigint,
	data jsonb,
	trx_id varchar(64),
	block_id varchar(64),
	fee bigint
);



create table if not exists public.transactions
(
	id varchar(64) not null,
	block_id varchar(64) not null,
	ref_block_num bigint,
	ref_block_prefix bigint,
	expiration timestamp,
	num_actions integer,
	delay_sec bigint,
	fee bigint,
	data jsonb,
	idx bigserial not null
		constraint transactions_pk
			primary key
);

create index if not exists transactions_block_id_index
	on public.transactions (block_id);

create index if not exists transactions_ref_block_num_index
	on public.transactions (ref_block_num);

create index if not exists transactions_expiration_index
	on public.transactions (expiration);

create index if not exists transactions_id_index
	on public.transactions (id);

create table if not exists public.actions
(
	id bigserial not null
		constraint actions_pkey
			primary key,
	account varchar(32),
	name varchar(32),
	data jsonb,
	trx_id varchar(64),
	block_id varchar(64),
	ref_block_num bigint,
	fee bigint,
	ref_block_prefix bigint
);

create unique index if not exists actions_id_uindex
	on public.actions (id);

create index if not exists actions_account_index
	on public.actions (account);

create index if not exists actions_name_index
	on public.actions (name);

create index if not exists actions_trx_id_index
	on public.actions (trx_id);

create index if not exists actions_block_id_index
	on public.actions (block_id);

create index if not exists actions_ref_block_num_index
	on public.actions (ref_block_num);

create table if not exists public.account_permissions
(
	id bigserial not null
		constraint account_permission_pkey
			primary key,
	account varchar(32) not null,
	permission varchar(64),
	pubkey varchar(256)
);



create unique index if not exists account_permission_id_uindex
	on public.account_permissions (id);

create index if not exists account_permission_account_index
	on public.account_permissions (account);

create index if not exists account_permission_pubkey_index
	on public.account_permissions (pubkey);

create table if not exists public.account_tokens
(
	name varchar(32) not null,
	chain varchar(16) not null,
	create_time timestamp default now(),
	update_time timestamp default now(),
	token_chain varchar(16) not null,
	symbol varchar(16) not null,
	amount bigint,
	constraint accounts_pk
		primary key (chain, name, token_chain, symbol)
);



