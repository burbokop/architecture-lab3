DROP TABLE IF EXISTS "virtual_mashins";
CREATE TABLE "virtual_mashins" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(64) NOT NULL UNIQUE,
  "cpuCount" INT NOT NULL,
  "discs" VARCHAR(256)
);

DROP TABLE IF EXISTS "discs";
CREATE TABLE "discs" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(50) NOT NULL,
  "diskSpace" bigint NOT NULL
);

INSERT INTO "discs" (id, name, diskSpace) VALUES (0, 'Intel', 4294967296);
INSERT INTO "discs" (id, name, diskSpace) VALUES (1, 'Toshiba', 17179869184);
INSERT INTO "discs" (id, name, diskSpace) VALUES (3, 'Toshiba2', 17179869184);

INSERT INTO "virtual_mashins" (name, cpuCount, discs) VALUES ('vm0', 4, "0, 1");
INSERT INTO "virtual_mashins" (name, cpuCount, discs) VALUES ('vm1', 4, "3");
