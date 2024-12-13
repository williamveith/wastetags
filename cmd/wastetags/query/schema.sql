CREATE TABLE IF NOT EXISTS mixtures (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chem_name TEXT,
    component_name TEXT,
    cas TEXT,
    percent TEXT,
    component_order INTEGER
);

CREATE TABLE IF NOT EXISTS chemicals (
    cas TEXT PRIMARY KEY,
    chem_name TEXT
);

CREATE TABLE IF NOT EXISTS locations (
    location TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS containers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    abbreviation TEXT,
    full_name TEXT GENERATED ALWAYS AS (
        CASE 
            WHEN abbreviation IS NOT NULL AND abbreviation != '' THEN abbreviation || ' ' || name
            ELSE name
        END
    ) VIRTUAL
);

CREATE TABLE IF NOT EXISTS units (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    abbreviation TEXT,
    full_name TEXT GENERATED ALWAYS AS (
        CASE 
            WHEN abbreviation IS NOT NULL AND abbreviation != '' THEN abbreviation || ' ' || name
            ELSE name
        END
    ) VIRTUAL
);

CREATE TABLE IF NOT EXISTS states (
    state TEXT PRIMARY KEY
);