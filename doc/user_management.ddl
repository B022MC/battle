-- User Management and Player-Dimension Records DDL
-- Assumed database: PostgreSQL

BEGIN;

-- 1) Ensure FK from game_account.user_id -> basic_user.id
ALTER TABLE IF EXISTS game_account
  ADD CONSTRAINT IF NOT EXISTS fk_game_account_user
  FOREIGN KEY (user_id) REFERENCES basic_user(id) ON UPDATE CASCADE ON DELETE RESTRICT;

-- 2) Restrict one store binding per game account
--    Using unique constraint on game_account_house.game_account_id so that a game account
--    can only be bound to a single house (store). If historical bindings are needed,
--    handle via archiving rows first.
ALTER TABLE IF EXISTS game_account_house
  ADD CONSTRAINT IF NOT EXISTS uk_game_account_house_account
  UNIQUE (game_account_id);

-- 3) Optional: also prevent duplicate (account, house) rows
ALTER TABLE IF EXISTS game_account_house
  ADD CONSTRAINT IF NOT EXISTS uk_game_account_house_account_house
  UNIQUE (game_account_id, house_gid);

-- 4) Enforce that a user can manage only one store at a time
--    Since game_shop_admin uses physical delete, use a partial unique index on (user_id)
--    for non-deleted rows. This allows historical rows with deleted_at IS NOT NULL.
CREATE UNIQUE INDEX IF NOT EXISTS uk_game_shop_admin_user_active
  ON game_shop_admin(user_id)
  WHERE deleted_at IS NULL;

-- 5) Create player-dimension record table for efficient querying
CREATE TABLE IF NOT EXISTS game_player_record (
  id               SERIAL PRIMARY KEY,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  battle_record_id INT4    NOT NULL,
  house_gid        INT4    NOT NULL,

  player_gid       INT8    NOT NULL,
  game_account_id  INT4,

  group_id         INT4    NOT NULL,
  room_uid         INT4    NOT NULL,
  kind_id          INT4    NOT NULL,
  base_score       INT4    NOT NULL,
  score_delta      INT4    NOT NULL DEFAULT 0,
  is_winner        BOOL    NOT NULL DEFAULT FALSE,
  battle_at        TIMESTAMPTZ NOT NULL,

  meta_json        JSONB   NOT NULL DEFAULT '{}'::jsonb
);

-- Add FKs and helpful indexes
ALTER TABLE IF EXISTS game_player_record
  ADD CONSTRAINT IF NOT EXISTS fk_gpr_battle
  FOREIGN KEY (battle_record_id) REFERENCES game_battle_record(id) ON UPDATE CASCADE ON DELETE CASCADE;

ALTER TABLE IF EXISTS game_player_record
  ADD CONSTRAINT IF NOT EXISTS fk_gpr_game_account
  FOREIGN KEY (game_account_id) REFERENCES game_account(id) ON UPDATE CASCADE ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_gpr_house_battleat ON game_player_record(house_gid, battle_at DESC);
CREATE INDEX IF NOT EXISTS idx_gpr_player_battleat ON game_player_record(player_gid, battle_at DESC);
CREATE INDEX IF NOT EXISTS idx_gpr_account_battleat ON game_player_record(game_account_id, battle_at DESC);
CREATE INDEX IF NOT EXISTS idx_gpr_group_room ON game_player_record(group_id, room_uid);

COMMIT;

