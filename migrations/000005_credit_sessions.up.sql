CREATE TABLE credit_sessions (
    id varchar(100) PRIMARY KEY,
    status integer,
    borrower varchar(42),
    account varchar(42),
    since integer,
    closed_at integer,
    initial_amount varchar(80),
    score double precision,
    credit_manager varchar(42),
    borrowed_amount varchar(80),
    profit varchar(80),
    profit_percent double precision,
    total_value varchar(80),
    health_factor integer,
    name text,
    background text
);
ALTER TABLE ONLY credit_sessions
    ADD CONSTRAINT credit_sessions_credit_manager_fkey FOREIGN KEY (credit_manager) REFERENCES credit_managers(address);

CREATE TABLE credit_session_snapshots (
    id SERIAL PRIMARY KEY,
    block_num bigint,
    session_id varchar(100),
    borrowed_amount_bi varchar(80),
    borrowed_amount double precision,
    total_value_bi varchar(80),
    total_value double precision,
    total_value_eth double precision,
    total_value_usd double precision,
    balances jsonb,
    cumulative_index character varying(80),
    health_factor integer,
    borrower character varying(42)
);


ALTER TABLE ONLY credit_session_snapshots
    ADD CONSTRAINT credit_session_snapshots_block_num_fkey FOREIGN KEY (block_num) REFERENCES blocks(id) ON DELETE CASCADE;
ALTER TABLE ONLY credit_session_snapshots
    ADD CONSTRAINT credit_session_snapshots_session_id_fkey FOREIGN KEY (session_id) REFERENCES credit_sessions(id);
