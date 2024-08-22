-- NOTE: This is just for resetting my database schema on my localhost
-- DO NOT USE THIS SCRIPT IN PRODUCTION

-- Drops all tables in the public schema
DO
$do$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public')
    LOOP
        EXECUTE 'DROP TABLE IF EXISTS "' || r.tablename || '" CASCADE;';
    END LOOP;
END
$do$;

-- Drops all functions in the information schema
DO
$$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT proname, pg_get_function_identity_arguments(p.oid) AS args
             FROM pg_proc p
             JOIN pg_namespace n ON p.pronamespace = n.oid
             WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
             AND pg_function_is_visible(p.oid)
             AND p.prokind = 'f'
    LOOP
        EXECUTE format('DROP FUNCTION IF EXISTS %I(%s);', r.proname, r.args);
    END LOOP;
END
$$;

-- Drops all procedures in the information schema
DO
$$
DECLARE
    r RECORD;
BEGIN
    FOR r IN SELECT proname, pg_get_function_identity_arguments(p.oid) AS args
             FROM pg_proc p
             JOIN pg_namespace n ON p.pronamespace = n.oid
             WHERE n.nspname NOT IN ('pg_catalog', 'information_schema')
             AND pg_function_is_visible(p.oid)
             AND p.prokind = 'p'
    LOOP
        EXECUTE format('DROP PROCEDURE IF EXISTS %I(%s);', r.proname, r.args);
    END LOOP;
END
$$;

--  Drops all views in the public schema
DO
$do$
DECLARE
    r RECORD;
BEGIN
    FOR r IN (SELECT viewname FROM pg_views WHERE schemaname = 'public')
    LOOP
        EXECUTE 'DROP VIEW IF EXISTS "' || r.viewname || '" CASCADE;';
    END LOOP;
END
$do$;

-- Drops all user-defined DOMAIN types in the current database
DO $$
DECLARE
    domain_rec RECORD;
BEGIN
    FOR domain_rec IN 
        SELECT domain_schema, domain_name
        FROM information_schema.domains
        WHERE domain_schema NOT IN ('pg_catalog', 'information_schema')
    LOOP
        EXECUTE format('DROP DOMAIN IF EXISTS %I.%I CASCADE;', 
                       domain_rec.domain_schema, domain_rec.domain_name);
        RAISE NOTICE 'Dropped domain: %.%', 
                     domain_rec.domain_schema, domain_rec.domain_name;
    END LOOP;
END $$;
