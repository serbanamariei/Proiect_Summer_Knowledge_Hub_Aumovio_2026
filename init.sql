-- asta i tabelu pentri link uri
CREATE TABLE IF NOT EXISTS siteuri_parsate(
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    linkuri JSONB
);

-- asta i tabelu pentru proxy uri
CREATE TABLE IF NOT EXISTS proxyuri(
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    address TEXT NOT NULL,
    username TEXT,
    password TEXT
);

-- initializez cu un proxy socks5
INSERT INTO proxyuri(type, address, username, password)
VALUES ('socks5','vpn-proxy:1080','admin','parola1109');