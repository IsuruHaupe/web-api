CREATE TABLE "contacts" (
  "id" bigserial PRIMARY KEY,
  "firstname" varchar NOT NULL,
  "lastname" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "home_address" varchar NOT NULL,
  "email" varchar NOT NULL,
  "phone_number" varchar NOT NULL
);

CREATE TABLE "skills" (
  "id" bigserial PRIMARY KEY,
  "skill_name" varchar NOT NULL,
  "skill_level" varchar NOT NULL,
  UNIQUE (skill_name, skill_level)
);

CREATE TABLE "contact_has_skill" (
  "contact_id" int NOT NULL,
  "skill_id" int NOT NULL,
  UNIQUE (contact_id, skill_id)
);

ALTER TABLE "contact_has_skill" ADD FOREIGN KEY ("contact_id") REFERENCES "contacts" ("id") ON DELETE CASCADE;

ALTER TABLE "contact_has_skill" ADD FOREIGN KEY ("skill_id") REFERENCES "skills" ("id") ON DELETE CASCADE;
