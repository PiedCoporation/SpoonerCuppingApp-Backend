BEGIN;

-- Enable required extension for UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Domain enums
DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'roastinglever_enum') THEN
		CREATE TYPE roastinglever_enum AS ENUM (
			'EXTRA_LIGHT',
			'LIGHT',
			'MEDIUM',
			'MEDIUM_DARK',
			'DARK'
		);
	END IF;
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'processing_enum') THEN
		CREATE TYPE processing_enum AS ENUM (
			'NATURAL',
			'HONEY',
			'SEMI',
			'WASHED',
			'ANAEROBIC'
		);
	END IF;
END $$;

-- Core tables
CREATE TABLE IF NOT EXISTS roles (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name text NOT NULL UNIQUE,
	description text NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS permissions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name text NOT NULL UNIQUE,
	description text NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	first_name text NULL,
	last_name text NULL,
	email text NOT NULL UNIQUE,
	phone text NULL,
	password text NOT NULL,
	is_verified boolean NOT NULL DEFAULT false,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	role_id uuid NOT NULL REFERENCES roles(id) ON DELETE RESTRICT
);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

CREATE TABLE IF NOT EXISTS refresh_tokens (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	token text NOT NULL,
	issued_at timestamptz NOT NULL,
	expires_at timestamptz NOT NULL,
	revoked boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);

CREATE TABLE IF NOT EXISTS events (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name text NOT NULL,
	date_of_event timestamptz NOT NULL,
	start_time timestamptz NOT NULL,
	end_time timestamptz NOT NULL,
	"limit" integer NOT NULL,
	total_current integer NOT NULL,
	number_samples integer NOT NULL,
	phone_contact text NOT NULL,
	email_contact text NOT NULL,
	invite_url text NULL,
	qr_image_url text NULL,
	is_public boolean NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);

CREATE TABLE IF NOT EXISTS event_addresses (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	province text NOT NULL,
	district text NOT NULL,
	longitude text NOT NULL,
	latitude text NOT NULL,
	ward text NOT NULL,
	street text NOT NULL,
	phone text NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_event_addresses_event_id ON event_addresses(event_id);

CREATE TABLE IF NOT EXISTS role_permissions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	permission_id uuid NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
	role_id uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
	UNIQUE (role_id, permission_id)
);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);

CREATE TABLE IF NOT EXISTS user_samples (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	name text NOT NULL,
	roasting_date timestamptz NOT NULL,
	roast_level roastinglever_enum NOT NULL,
	altitude_grow text NOT NULL,
	roastery_name text NOT NULL,
	roastery_address text NOT NULL,
	breed_name text NOT NULL,
	pre_processing processing_enum NOT NULL,
	grow_nation text NOT NULL,
	grow_address text NOT NULL,
	price double precision NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_user_samples_user_id ON user_samples(user_id);

CREATE TABLE IF NOT EXISTS event_samples (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	price text NOT NULL,
	rating integer NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	user_sample_id uuid NOT NULL REFERENCES user_samples(id) ON DELETE CASCADE,
	event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_event_samples_user_sample_id ON event_samples(user_sample_id);
CREATE INDEX IF NOT EXISTS idx_event_samples_event_id ON event_samples(event_id);

CREATE TABLE IF NOT EXISTS event_users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	is_accepted boolean NOT NULL,
	is_invited boolean NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_event_users_user_id ON event_users(user_id);
CREATE INDEX IF NOT EXISTS idx_event_users_event_id ON event_users(event_id);

CREATE TABLE IF NOT EXISTS user_tastings (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	rating integer NOT NULL CHECK (rating >= 1 AND rating <= 5),
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	event_sample_id uuid NOT NULL REFERENCES event_samples(id) ON DELETE CASCADE,
	event_user_id uuid NOT NULL REFERENCES event_users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_user_tastings_event_sample_id ON user_tastings(event_sample_id);
CREATE INDEX IF NOT EXISTS idx_user_tastings_event_user_id ON user_tastings(event_user_id);

CREATE TABLE IF NOT EXISTS tasting_notes (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	parent_name text NULL,
	child_name text NULL,
	grand_child_name text NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	user_tasting_id uuid NOT NULL REFERENCES user_tastings(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_tasting_notes_user_tasting_id ON tasting_notes(user_tasting_id);

COMMIT;
