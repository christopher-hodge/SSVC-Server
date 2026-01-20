CREATE TABLE players (
    player_id      UUID PRIMARY KEY,
    display_name   TEXT NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login_at  TIMESTAMP
);

CREATE TABLE item_definitions (
    item_type_id   TEXT PRIMARY KEY,
    name           TEXT NOT NULL,
    rarity         TEXT NOT NULL,
    max_stack      INT NOT NULL DEFAULT 1,
    base_stats     JSONB NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE loot_tables (
    loot_table_id  TEXT PRIMARY KEY,
    rolls          INT NOT NULL
);

CREATE TABLE loot_table_entries (
    loot_table_id  TEXT REFERENCES loot_tables(loot_table_id),
    item_type_id   TEXT REFERENCES item_definitions(item_type_id),
    weight         INT NOT NULL,
    PRIMARY KEY (loot_table_id, item_type_id)
);

CREATE TABLE item_instances (
    instance_id    UUID PRIMARY KEY,
    item_type_id   TEXT NOT NULL REFERENCES item_definitions(item_type_id),

    owner_type     TEXT NOT NULL CHECK (
        owner_type IN ('player', 'world', 'escrow', 'destroyed')
    ),

    owner_id       UUID,
    state          TEXT NOT NULL CHECK (
        state IN ('inventory', 'world', 'trade', 'destroyed')
    ),

    rolled_stats   JSONB NOT NULL,
    version        INT NOT NULL DEFAULT 1,

    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE world_items (
    instance_id    UUID PRIMARY KEY REFERENCES item_instances(instance_id),
    map_id         TEXT NOT NULL,
    pos_x          FLOAT NOT NULL,
    pos_y          FLOAT NOT NULL,
    pos_z          FLOAT NOT NULL,
    spawned_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE VIEW player_inventory AS
SELECT
    instance_id,
    item_type_id,
    rolled_stats,
    owner_id AS player_id
FROM item_instances
WHERE owner_type = 'player'
  AND state = 'inventory';

CREATE TABLE trades (
    trade_id       UUID PRIMARY KEY,
    initiator_id   UUID NOT NULL REFERENCES players(player_id),
    recipient_id   UUID NOT NULL REFERENCES players(player_id),

    status         TEXT NOT NULL CHECK (
        status IN ('pending', 'accepted', 'cancelled', 'completed')
    ),

    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE item_destruction_log (
    instance_id    UUID PRIMARY KEY,
    reason         TEXT NOT NULL,
    destroyed_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE server_events (
    event_id       BIGSERIAL PRIMARY KEY,
    player_id      UUID,
    event_type     TEXT NOT NULL,
    payload        JSONB,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);