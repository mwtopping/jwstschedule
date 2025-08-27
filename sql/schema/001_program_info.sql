CREATE TABLE program_info (
	id INTEGER PRIMARY KEY,
	created_at INTEGER NOT NULL,
	updated_at INTEGER NOT NULL,
	title TEXT NOT NULL,
	pi TEXT NOT NULL,   
	eap INTEGER NOT NULL,
	primetime REAL NOT NULL,      
	paralleltime REAL NOT NULL,   
	InstrumentMode TEXT NOT NULL,
	ProgramType TEXT NOT NULL
);

