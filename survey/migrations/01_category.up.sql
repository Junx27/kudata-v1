CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL
);


INSERT INTO categories (name) VALUES
('Teknologi'),
('Kesehatan'),
('Pendidikan'),
('Keuangan'),
('Gaya Hidup'),
('Wisata'),
('Kuliner'),
('Olahraga'),
('Seni & Budaya'),
('Lingkungan');
