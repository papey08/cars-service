CREATE TABLE "cars" (
    "id" SERIAL PRIMARY KEY,
    "reg_num" VARCHAR(9) UNIQUE,
    "model_id" INTEGER,
    "year" INTEGER,
    "owner_id" INTEGER
);

CREATE TABLE "models" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100),
    "mark_id" INTEGER
);

CREATE TABLE "marks" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100)
);

CREATE TABLE "owners" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100),
    "surname" VARCHAR(100),
    "patronymic" VARCHAR(100)
);
