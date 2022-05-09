for DB in $(psql -U user -t -c "SELECT datname FROM pg_database WHERE datname NOT IN ('postgres', 'template0', 'template1')"); do
  psql -U user -d $DB -c 'CREATE EXTENSION IF NOT EXISTS "uuid-ossp"' 
  psql -U user -d $DB -c 'CREATE EXTENSION IF NOT EXISTS "postgis"' 
done
