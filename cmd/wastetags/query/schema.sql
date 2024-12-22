PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS chemicals (
    cas TEXT PRIMARY KEY,
    chem_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS mixtures (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chem_name TEXT NOT NULL,
    component_name TEXT NOT NULL,
    cas TEXT NOT NULL,
    percent TEXT NOT NULL,
    FOREIGN KEY (chem_name) REFERENCES alias (display_name) ON DELETE NO ACTION,
    FOREIGN KEY (component_name) REFERENCES alias (display_name) ON DELETE NO ACTION,
    FOREIGN KEY (cas) REFERENCES chemicals (cas) ON DELETE NO ACTION,
    UNIQUE (chem_name, cas)
);

CREATE TABLE IF NOT EXISTS alias (
    display_name TEXT PRIMARY KEY, -- Name printed on wastetag
    internal_name TEXT NOT NULL -- Name registered in EHS database
);

CREATE TABLE IF NOT EXISTS locations (
    location TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS containers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    abbreviation TEXT,
    full_name TEXT GENERATED ALWAYS AS (
        CASE 
            WHEN abbreviation IS NOT NULL AND abbreviation != '' THEN abbreviation || ' ' || name
            ELSE name
        END
    ) VIRTUAL,
    UNIQUE (name, abbreviation)
);

CREATE TABLE IF NOT EXISTS units (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    abbreviation TEXT,
    full_name TEXT GENERATED ALWAYS AS (
        CASE 
            WHEN abbreviation IS NOT NULL AND abbreviation != '' THEN abbreviation || ' ' || name
            ELSE name
        END
    ) VIRTUAL,
    UNIQUE (name, abbreviation)
);

CREATE TABLE IF NOT EXISTS states (
    state TEXT PRIMARY KEY
);

CREATE VIEW IF NOT EXISTS mixtures_internal_names AS
SELECT 
    m.id,
    m.chem_name        AS chem_display_name,
    a1.internal_name   AS chem_internal_name,
    m.component_name   AS component_display_name,
    a2.internal_name   AS component_internal_name,
    m.cas,
    m.percent
FROM mixtures m
JOIN alias a1 
    ON m.chem_name = a1.display_name
JOIN alias a2 
    ON m.component_name = a2.display_name;