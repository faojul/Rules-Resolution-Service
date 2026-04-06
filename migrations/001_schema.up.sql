CREATE TABLE steps (
    key TEXT PRIMARY KEY,
    name TEXT,
    description TEXT,
    position INT
);

CREATE TABLE defaults (
    step_key TEXT,
    trait_key TEXT,
    value JSONB,
    PRIMARY KEY (step_key, trait_key)
);

CREATE TABLE overrides (
    id TEXT PRIMARY KEY,
    step_key TEXT,
    trait_key TEXT,

    selector JSONB,

    value JSONB,
    specificity INT,

    effective_date TIMESTAMP,
    expires_date TIMESTAMP,

    status TEXT,
    description TEXT,

    created_by TEXT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);