-- Пользователь только на чтение для реплик Python backend (store_python)
CREATE USER shop_readonly WITH PASSWORD 'readonly';
GRANT CONNECT ON DATABASE shopapi TO shop_readonly;
\c shopapi
GRANT USAGE ON SCHEMA public TO shop_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO shop_readonly;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres IN SCHEMA public GRANT SELECT ON TABLES TO shop_readonly;
