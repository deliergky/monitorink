CREATE TABLE response (
    url character varying(256) NOT NULL,
    status_code integer NOT NULL,
    response_time integer NOT NULL,
    matched boolean,
    created_at timestamp with time zone DEFAULT timezone('UTC'::text, now()) NOT NULL
);
