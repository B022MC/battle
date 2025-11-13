#!/bin/bash
#
# This script generates a full, clean DDL schema for the battle-tiles database.
# It uses pg_dump to extract the schema without any data, ownership, or privileges.
#
# Requirements:
# - `pg_dump` must be installed and in the system's PATH.
# - The user running the script must have read access to the database.
#
# Usage:
# 1. Replace placeholder values for DB_USER, DB_HOST, DB_PORT, and DB_NAME.
# 2. Run the script from the repository root:
#    bash doc/generate_full_schema.sh > doc/full_schema.ddl
#

# --- Configuration ---
DB_USER="your_db_user"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="battle_tiles_db"

# --- Dump Command ---
# --schema-only: Dumps only the object definitions (schema), not data.
# --no-owner: Prevents dumping of object ownership commands.
# --no-privileges: Prevents dumping of access privileges (GRANT/REVOKE).
# --clean: Adds DROP commands to clean the database before recreating.
# -U: Database user
# -h: Database host
# -p: Database port
# -d: Database name

PGPASSWORD="your_db_password" pg_dump \
  --schema-only \
  --no-owner \
  --no-privileges \
  --clean \
  -U "${DB_USER}" \
  -h "${DB_HOST}" \
  -p "${DB_PORT}" \
  -d "${DB_NAME}"

echo "Dump command executed. Redirect output to doc/full_schema.ddl to save the file."

