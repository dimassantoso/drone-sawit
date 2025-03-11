-- estates table
CREATE TABLE IF NOT EXISTS estates (
                                       "id"     varchar(36) PRIMARY KEY,
    "width"  INT NOT NULL CHECK (width > 0 AND width <= 50000),
    "length" INT NOT NULL CHECK (length > 0 AND length <= 50000),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    );

-- estate_tress table
CREATE TABLE IF NOT EXISTS estate_trees
(
    id        varchar(36) PRIMARY KEY,
    estate_id varchar(36) REFERENCES estates (id) ON DELETE CASCADE,
    x         INT NOT NULL CHECK (x > 0),
    y         INT NOT NULL CHECK (y > 0),
    height    INT NOT NULL CHECK (height > 0 AND height < 30),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_estate_id ON estate_trees USING btree (estate_id);
CREATE INDEX idx_tree_coords ON estate_trees USING btree (estate_id, x, y);