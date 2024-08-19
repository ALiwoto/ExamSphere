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
    LOOP
        EXECUTE format('DROP FUNCTION IF EXISTS %I(%s);', r.proname, r.args);
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