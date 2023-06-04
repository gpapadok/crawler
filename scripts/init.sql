create table links (
    parent text,
    url text primary key,
    crawled_at timestamp with time zone default clock_timestamp()
);
