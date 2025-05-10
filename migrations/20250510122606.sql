-- Modify "choices" table
ALTER TABLE "choices" ADD COLUMN "uid" character varying NOT NULL;
-- Create index "choices_uid_key" to table: "choices"
CREATE UNIQUE INDEX "choices_uid_key" ON "choices" ("uid");
-- Modify "questions" table
ALTER TABLE "questions" ADD COLUMN "uid" character varying NOT NULL;
-- Create index "questions_uid_key" to table: "questions"
CREATE UNIQUE INDEX "questions_uid_key" ON "questions" ("uid");
