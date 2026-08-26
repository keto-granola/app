CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  description TEXT NOT NULL,
  ingredients JSONB NOT NULL,
  nutrition JSONB NOT NULL,
  weight_g INTEGER NOT NULL,
  dietary_tags TEXT[] NOT NULL DEFAULT '{}',
  allergens TEXT[] NOT NULL DEFAULT '{}',
  price_cents INTEGER NOT NULL,
  currency CHAR(3) NOT NULL DEFAULT 'AUD',
  image_storage_key TEXT NOT NULL,
  image_alt TEXT NOT NULL
);


CREATE TYPE order_status AS ENUM ('confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded');

CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  order_number TEXT UNIQUE NOT NULL,
  status order_status NOT NULL DEFAULT 'confirmed',

  -- Customer info
  customer_name TEXT NOT NULL,
  customer_email TEXT NOT NULL,
  customer_phone TEXT,
  shipping_line1 TEXT NOT NULL,
  shipping_line2 TEXT,
  shipping_city TEXT NOT NULL,
  shipping_state TEXT,
  shipping_postcode TEXT NOT NULL,
  shipping_country TEXT NOT NULL DEFAULT 'AU',

  -- Stripe
  stripe_payment_intent_id TEXT UNIQUE NOT NULL,
  stripe_payment_status TEXT NOT NULL,       
  amount_total_cents INTEGER NOT NULL,
  currency CHAR(3) NOT NULL DEFAULT 'AUD',

  -- Shipping
  shipping_method             TEXT,
  tracking_number             TEXT,
  tracking_url                TEXT,
  estimated_delivery_date     DATE,
  shipped_at                  TIMESTAMPTZ,
  delivered_at                TIMESTAMPTZ
);

CREATE TABLE order_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

  -- Item
  product_id UUID NOT NULL REFERENCES products(id),
  product_name TEXT NOT NULL,
  quantity INTEGER NOT NULL CHECK (quantity > 0),
  unit_price_cents INTEGER NOT NULL,
  total_price_cents INTEGER NOT NULL
);

CREATE TABLE inventory (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  restock_threshold INTEGER NOT NULL CHECK (restock_threshold > 0)
);