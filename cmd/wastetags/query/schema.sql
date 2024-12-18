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
    component_order INTEGER,
    FOREIGN KEY (cas) REFERENCES chemicals (cas) ON DELETE CASCADE
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
    ) VIRTUAL
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
    ) VIRTUAL
);

CREATE TABLE IF NOT EXISTS states (
    state TEXT PRIMARY KEY
);