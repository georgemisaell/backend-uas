CREATE TABLE students (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID REFERENCES users(id),
	student_id VARCHAR(20) UNIQUE NOT NULL,
	program_study VARCHAR(100),
	academy_year VARCHAR(10),
	advisor_id UUID REFERENCES lecturers(id),
	created_at TIMESTAMP DEFAULT NOW()
	updated_at TIMESTAMP DEFAULT NOW()
);