-- Enable foreign key constraints to enforce relationships between tables
PRAGMA foreign_keys = ON;
PRAGMA journal_mode=WAL;

-- Create a table to store chemicals with a unique CAS identifier and chemical name
CREATE TABLE IF NOT EXISTS chemicals (
    cas TEXT PRIMARY KEY, -- Unique CAS number for the chemical
    chem_name TEXT NOT NULL -- Name of the chemical
);

-- Create a table to store mixtures of chemicals, referencing alias and chemicals tables
CREATE TABLE IF NOT EXISTS mixtures (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Unique ID for the mixture
    chem_name TEXT NOT NULL, -- Name of the chemical in the mixture
    component_name TEXT NOT NULL, -- Name of the component in the mixture
    cas TEXT NOT NULL, -- CAS number of the chemical
    percent TEXT NOT NULL, -- Percentage of the component in the mixture
    FOREIGN KEY (chem_name) REFERENCES alias (display_name) ON DELETE NO ACTION, -- Enforces relationship to alias table
    FOREIGN KEY (component_name) REFERENCES alias (display_name) ON DELETE NO ACTION, -- Enforces relationship to alias table
    FOREIGN KEY (cas) REFERENCES chemicals (cas) ON DELETE NO ACTION, -- Enforces relationship to chemicals table
    UNIQUE (chem_name, cas) -- Ensures unique pairing of chemical name and CAS number
);

-- Create a table to store aliases for chemicals with display and internal names
CREATE TABLE IF NOT EXISTS alias (
    display_name TEXT PRIMARY KEY, -- Name printed on wastetag
    internal_name TEXT NOT NULL -- Name registered in EHS database
);

-- Create a table to store unique locations
CREATE TABLE IF NOT EXISTS locations (
    location TEXT PRIMARY KEY -- Unique identifier for each location
);

-- Create a table to store container types with optional abbreviations
CREATE TABLE IF NOT EXISTS containers (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Unique ID for the container
    name TEXT NOT NULL, -- Name of the container
    abbreviation TEXT, -- Abbreviation of the container name
    full_name TEXT GENERATED ALWAYS AS (
        CASE 
            WHEN abbreviation IS NOT NULL AND abbreviation != '' THEN abbreviation || ' ' || name
            ELSE name
        END
    ) VIRTUAL, -- Virtual column (not stored in the table) Auto-generated full name based on abbreviation and name
    UNIQUE (name, abbreviation) -- Ensures unique combination of name and abbreviation
);

-- Create a table to store units with optional abbreviations
CREATE TABLE IF NOT EXISTS units (
    id INTEGER PRIMARY KEY AUTOINCREMENT, -- Unique ID for the unit
    name TEXT NOT NULL, -- Name of the unit
    abbreviation TEXT, -- Abbreviation of the unit name
    full_name TEXT GENERATED ALWAYS AS (
        CASE 
            WHEN abbreviation IS NOT NULL AND abbreviation != '' THEN abbreviation || ' ' || name
            ELSE name
        END
    ) VIRTUAL, -- Virtual column (not stored in the table). Auto-generated full name based on abbreviation and name
    UNIQUE (name, abbreviation) -- Ensures unique combination of name and abbreviation
);

-- Create a table to store physical states
CREATE TABLE IF NOT EXISTS states (
    state TEXT PRIMARY KEY -- Unique identifier for the physical state
);

-- Create a view to link mixtures with their internal names for easier querying
CREATE VIEW IF NOT EXISTS mixtures_internal_names AS
SELECT 
    m.id, -- ID of the mixture
    m.chem_name        AS chem_display_name, -- Display name of the chemical
    a1.internal_name   AS chem_internal_name, -- Internal name of the chemical
    m.component_name   AS component_display_name, -- Display name of the component
    a2.internal_name   AS component_internal_name, -- Internal name of the component
    m.cas, -- CAS number of the chemical
    m.percent -- Percentage of the component in the mixture
FROM mixtures m
JOIN alias a1 
    ON m.chem_name = a1.display_name -- Join with alias table for chemical names
JOIN alias a2 
    ON m.component_name = a2.display_name; -- Join with alias table for component names

-- Trigger to update chem_name in mixtures when display_name in alias is updated
CREATE TRIGGER IF NOT EXISTS update_chem_name_on_alias_update
AFTER UPDATE OF display_name ON alias -- Trigger fires after updating display_name in alias
FOR EACH ROW
BEGIN
    UPDATE mixtures
    SET chem_name = NEW.display_name -- Update chem_name to the new display_name
    WHERE chem_name = OLD.display_name; -- Match the old display_name
END;

-- Trigger to update component_name in mixtures when display_name in alias is updated
CREATE TRIGGER IF NOT EXISTS update_component_name_on_alias_update
AFTER UPDATE OF display_name ON alias -- Trigger fires after updating display_name in alias
FOR EACH ROW
BEGIN
    UPDATE mixtures
    SET component_name = NEW.display_name -- Update component_name to the new display_name
    WHERE component_name = OLD.display_name; -- Match the old display_name
END;
