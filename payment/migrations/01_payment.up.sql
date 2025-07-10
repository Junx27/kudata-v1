CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  survey_id INT NOT NULL,
  status BOOLEAN NOT NULL
);
