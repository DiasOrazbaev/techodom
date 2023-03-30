create table if not exists links
(
    id           serial primary key,
    active_link  varchar(255) not null,
    history_link varchar(255) not null
);

create unique index if not exists links_active_link_history_link_uindex
    on public.links (active_link, history_link);