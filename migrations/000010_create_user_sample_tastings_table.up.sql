-- Create user_sample_tastings table
CREATE TABLE IF NOT EXISTS user_sample_tastings (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	parent_name text NULL,
	child_name text NULL,
	grand_child_name text NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	user_sample_id uuid NOT NULL REFERENCES user_samples(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_user_sample_tastings_user_sample_id ON user_sample_tastings(user_sample_id);
