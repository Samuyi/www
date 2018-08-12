CREATE TABLE IF NOT EXISTS users (
id  uuid DEFAULT uuid_generate_v4() UNIQUE,
display_name text NOT NULL UNIQUE,
first_name text NOT NULL,
last_name text NOT NULL,
ratings INTEGER DEFAULT 0,
password text NOT NULL,
email text NOT NULL UNIQUE,
avatar text,
active BOOLEAN DEFAULT FALSE,
created_at TIMESTAMP DEFAULT NOW(),
updated_at TIMESTAMP,
PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS locations(
    location_id uuid DEFAULT uuid_generate_v4() UNIQUE,
    city text NOT NULL UNIQUE,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    state text NOT NULL,
    country text NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    PRIMARY KEY(location_id)
);

CREATE TABLE IF NOT EXISTS items (
    id  uuid DEFAULT uuid_generate_v4() UNIQUE,
    name text NOT NULL,
    closed BOOLEAN DEFAULT FALSE,
    city text REFERENCES locations(city),
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    instruction VARCHAR, 
    phone_no VARCHAR(14) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (id)
);