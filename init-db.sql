-- Minimal init: create dev/test databases only. Run migrations separately.

-- Ensure dblink extension is available for cross-database creation
CREATE EXTENSION IF NOT EXISTS dblink;

-- Ensure databases exist
DO
$$
BEGIN
	IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'hackathon_development') THEN
		PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE hackathon_development');
	END IF;
	IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'hackathon_test') THEN
		PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE hackathon_test');
	END IF;
END
$$ LANGUAGE plpgsql;

-- No seeding here; use `make db-seed-sql` after migrations.