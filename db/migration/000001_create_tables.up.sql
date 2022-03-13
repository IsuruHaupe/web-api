CREATE TABLE "contacts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "firstname" varchar NOT NULL,
  "lastname" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "home_address" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone_number" varchar NOT NULL
);

CREATE TABLE "skills" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "skill_name" varchar NOT NULL,
  "skill_level" varchar NOT NULL, 
  UNIQUE ("owner", "skill_name", "skill_level")
);

CREATE TABLE "contact_has_skill" (
  "owner" varchar NOT NULL,
  "contact_id" int NOT NULL,
  "skill_id" int NOT NULL, 
  UNIQUE ("owner", "contact_id", "skill_id")
);

CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_last_changed" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "create_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "session_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "contacts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "skills" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "contact_has_skill" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;

ALTER TABLE "contact_has_skill" ADD FOREIGN KEY ("contact_id") REFERENCES "contacts" ("id") ON DELETE CASCADE;

ALTER TABLE "contact_has_skill" ADD FOREIGN KEY ("skill_id") REFERENCES "skills" ("id") ON DELETE CASCADE;
