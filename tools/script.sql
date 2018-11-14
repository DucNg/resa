DROP TABLE Voucher;
DROP TABLE Session;
DROP TABLE Invite;
DROP TABLE Administrateur;

CREATE TABLE Invite (
	id_invite INTEGER PRIMARY KEY,
	nom TEXT,
	prenom TEXT,
	mail TEXT NOT NULL,
	mdp TEXT NOT NULL,
	numtel TEXT,
	parrain INTEGER REFERENCES id_invite
);

CREATE TABLE Voucher (
	id_voucher INTEGER PRIMARY KEY,
	code TEXT,
	expiration TIMESTAMP,
	proprietaire INTEGER,
	FOREIGN KEY (proprietaire) REFERENCES Invite(id_invite)
);

CREATE TABLE Administrateur (
	id_admin INTEGER PRIMARY KEY,
	login TEXT,
	mdp TEXT
);

CREATE TABLE Session (
	token TEXT NOT NULL PRIMARY KEY,
	id_user INTEGER REFERENCES Invite(id_invite)
);

CREATE TABLE AdminSession (
	token TEXT NOT NULL PRIMARY KEY,
	id_user INTEGER REFERENCES Invite(id_admin)
)