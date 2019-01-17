-- Database: Trader

-- DROP DATABASE "Trader";

CREATE DATABASE "Trader"
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'English_United States.1252'
    LC_CTYPE = 'English_United States.1252'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;
 
\c "Trader"

CREATE TABLE Contact (
	contactid SERIAL PRIMARY KEY,
	clientid INT NOT NULL, 
	broker VARCHAR(100) NOT NULL,
	country VARCHAR(100) NOT NULL,
	language VARCHAR(100) NOT NULL,
	identifier VARCHAR(100) NOT NULL,
    CONSTRAINT u_client UNIQUE (clientid, language, identifier)
);
 
CREATE TABLE SubscriptionType (
	subscriptionTypeid SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	CONSTRAINT u_account_name UNIQUE (name)
);

CREATE TABLE Subscription (
	subscriptionid SERIAL PRIMARY KEY,
	subscriptionTypeid INT NOT NULL,
	account VARCHAR(100) NULL,
	CONSTRAINT u_subs_account UNIQUE (subscriptionTypeid, account),
	FOREIGN KEY (subscriptionTypeid) REFERENCES SubscriptionType (subscriptionTypeid)
);

CREATE TABLE ContactSubscription (
	contactid INT NOT NULL,
	subscriptionid INT NOT NULL,
	PRIMARY KEY (contactid, subscriptionid),
	FOREIGN KEY (contactid) REFERENCES Contact (contactid),
	FOREIGN KEY (subscriptionid) REFERENCES Subscription (subscriptionid)
);

CREATE USER trader1 WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE "Trader" TO trader1;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO trader1;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO trader1;