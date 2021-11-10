CREATE TABLE IF NOT EXISTS Taxonomy (
	ID INT PRIMARY KEY,
	Kingdom TEXT,
	Phylum TEXT,
	Class TEXT,
	Orders TEXT,
	Family TEXT,
	Genus TEXT,
	Species TEXT,
	Citation TEXT,
	DB TEXT
);

CREATE TABLE IF NOT EXISTS Common (
	ID INT,
	Name TEXT,
	CONSTRAINT fk_taxonomy_common FOREIGN KEY (ID) REFERENCES Taxonomy(ID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IX_taxonomy_id ON Taxonomy (ID);
CREATE INDEX IX_common_id ON Common (ID);
