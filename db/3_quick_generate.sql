CREATE TABLE image_generations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    style_id UUID REFERENCES styles(id) ON DELETE SET NULL,
    room_type VARCHAR DEFAULT '',
    prompt TEXT NOT NULL,
    image_url TEXT NOT NULL,
    generated_image_url TEXT NOT NULL,
    status VARCHAR NOT NULL DEFAULT 'pending',
    error_message TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);