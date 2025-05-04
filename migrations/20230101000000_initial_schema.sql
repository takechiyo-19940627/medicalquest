-- Create Questions table
CREATE TABLE "questions" (
    "id" serial PRIMARY KEY,
    "reference_code" varchar NULL,
    "title" varchar NOT NULL,
    "content" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create Choices table
CREATE TABLE "choices" (
    "id" serial PRIMARY KEY,
    "question_id" integer NOT NULL,
    "content" varchar NOT NULL,
    "is_correct" boolean NOT NULL DEFAULT false,
    CONSTRAINT "choices_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE CASCADE
);