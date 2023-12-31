create table contracts
(
    id              bigserial                 not null
    constraint contracts_pk
    primary key,
    created         timestamptz default now() not null,
    updated         timestamptz default now() not null,
    provider_id     bigint                    not null references providers (id),
    delegate_pubkey text                      not null check ( delegate_pubkey != '' ),
    client_pubkey   text                      not null check ( client_pubkey != '' ),
    height          bigint                    not null check ( height > 0 ),
    contract_type   text                      not null references contract_types (val),
    duration        bigint                    not null,
    rate_asset      text                      not null,
    rate_amount     bigint                    not null,
    auth            text                      not null references auth_types (val),
    open_cost       bigint                    not null,
    deposit         bigint                    not null,
    queries_per_minute      bigint                    not null,
    settlement_duration     bigint                    not null,
    nonce           bigint                    not null DEFAULT 0,
    paid            bigint                    not null DEFAULT 0,
    reserve_contrib_asset   bigint                    not null DEFAULT 0,
    reserve_contrib_usd     bigint                    not null DEFAULT 0
);

alter table contracts
    add constraint pubkey_prov_dlgt_uniq unique (provider_id, delegate_pubkey);

---- create above / drop below ----

drop table contracts;
