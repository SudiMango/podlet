-- +goose Up
-- +goose StatementBegin
ALTER TABLE dns_records
    DROP CONSTRAINT dns_records_server_id_fkey,
    ADD CONSTRAINT dns_records_server_id_fkey
        FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE;

ALTER TABLE dns_records
    DROP CONSTRAINT dns_records_tunnel_id_fkey,
    ADD CONSTRAINT dns_records_tunnel_id_fkey
        FOREIGN KEY (tunnel_id) REFERENCES tunnels (id) ON DELETE CASCADE;

ALTER TABLE tunnels
    DROP CONSTRAINT tunnels_server_id_fkey,
    ADD CONSTRAINT tunnels_server_id_fkey
        FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE;

ALTER TABLE server_agent_jobs
    DROP CONSTRAINT server_agent_jobs_server_id_fkey,
    ADD CONSTRAINT server_agent_jobs_server_id_fkey
        FOREIGN KEY (server_id) REFERENCES servers (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE server_agent_jobs
    DROP CONSTRAINT server_agent_jobs_server_id_fkey,
    ADD CONSTRAINT server_agent_jobs_server_id_fkey
        FOREIGN KEY (server_id) REFERENCES servers (id);

ALTER TABLE tunnels
    DROP CONSTRAINT tunnels_server_id_fkey,
    ADD CONSTRAINT tunnels_server_id_fkey
        FOREIGN KEY (server_id) REFERENCES servers (id);

ALTER TABLE dns_records
    DROP CONSTRAINT dns_records_tunnel_id_fkey,
    ADD CONSTRAINT dns_records_tunnel_id_fkey
        FOREIGN KEY (tunnel_id) REFERENCES tunnels (id);

ALTER TABLE dns_records
    DROP CONSTRAINT dns_records_server_id_fkey,
    ADD CONSTRAINT dns_records_server_id_fkey
        FOREIGN KEY (server_id) REFERENCES servers (id);
-- +goose StatementEnd
