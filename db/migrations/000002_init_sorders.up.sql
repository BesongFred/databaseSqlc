CREATE TABLE orders(
  "id" VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid()::varchar(36),
  "amount" VARCHAR(36) NOT NULL,
  "number" VARCHAR(36) NOT NULL,
  "status" TEXT DEFAULT 'PENDING',
  "created_at" TIMESTAMP DEFAULT now()
);