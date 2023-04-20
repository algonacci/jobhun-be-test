CREATE TABLE Mahasiswa (
  Id INT PRIMARY KEY,
  Nama VARCHAR(255),
  Usia INT,
  Gender INT,
  Tanggal_Registrasi DATETIME
);

CREATE TABLE Jurusan (
  Id INT PRIMARY KEY,
  Nama_Jurusan VARCHAR(255)
);

CREATE TABLE Hobi (
  Id INT PRIMARY KEY,
  Nama_Hobi VARCHAR(255)
);

CREATE TABLE Mahasiswa_Hobi (
  Id_Mahasiswa INT,
  Id_Hobi INT,
  PRIMARY KEY (Id_Mahasiswa, Id_Hobi)
);
