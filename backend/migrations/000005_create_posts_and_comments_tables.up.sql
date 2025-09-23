-- posts
CREATE TABLE IF NOT EXISTS posts (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title varchar(255) NOT NULL,
    content text NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
    event_id uuid NULL REFERENCES events(id) ON DELETE SET NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_posts_event_id ON posts(event_id);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);

-- post_images
CREATE TABLE IF NOT EXISTS post_images (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    url text NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    UNIQUE(url)
);
CREATE INDEX IF NOT EXISTS idx_post_images_post_id ON post_images(post_id);

-- post_likes
CREATE TABLE IF NOT EXISTS post_likes (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (post_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_post_likes_post_id ON post_likes(post_id);
CREATE INDEX IF NOT EXISTS idx_post_likes_user_id ON post_likes(user_id);

-- post_comments
CREATE TABLE IF NOT EXISTS post_comments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    content text NOT NULL,
	is_deleted boolean NOT NULL DEFAULT false,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id uuid NULL REFERENCES post_comments(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_post_comments_post_id ON post_comments(post_id);
CREATE INDEX IF NOT EXISTS idx_post_comments_parent_id ON post_comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_post_comments_user_id ON post_comments(user_id);