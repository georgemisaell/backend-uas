# Sistem Pelaporan Prestasi Mahasiswa (Backend API)

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Fiber Framework](https://img.shields.io/badge/Fiber-v2-black?style=flat&logo=gofiber)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)
![MongoDB](https://img.shields.io/badge/MongoDB-6.0+-47A248?style=flat&logo=mongodb)

Backend service untuk sistem pelaporan dan validasi prestasi mahasiswa. Dibangun menggunakan arsitektur **Hybrid Database** (PostgreSQL untuk relasional & RBAC, MongoDB untuk data prestasi dinamis).

## 📋 Fitur Utama

- **Autentikasi JWT** (Login, Refresh Token).
- **Role-Based Access Control (RBAC)**: Admin, Mahasiswa, Dosen Wali.
- **Manajemen Prestasi**:
  - **Hybrid Storage:** Data referensi & status di PostgreSQL, detail dinamis di MongoDB.
  - **Workflow Status:** _Draft_ -> _Submitted_ -> _Verified_ / _Rejected_.
  - **Validasi Ketat:** Hak akses Dosen Wali terhadap Mahasiswa bimbingan.
- **Manajemen User & Mahasiswa**.

## 🛠️ Tech Stack

- **Language:** Golang
- **Framework:** Fiber v2
- **Database:** PostgreSQL (Relational) & MongoDB (NoSQL)
- **Migration Tool:** Golang Migrate

## 🚀 Cara Menjalankan Project

### 1. Prasyarat

Pastikan tool berikut sudah terinstall di komputer Anda:

- [Go](https://go.dev/dl/)
- [PostgreSQL](https://www.postgresql.org/)
- [MongoDB](https://www.mongodb.com/)
- [Golang Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### 2. Konfigurasi Environment

Buat file `.env` di root folder project dan isi konfigurasi berikut:

```env
APP_PORT=3000
API_KEY={isi bebas}
POSGRES_URI=postgres://postgres:root@localhost:5432/uas?sslmode=disable
MONGO_URI=mongodb://localhost:27017/uas
JWT_SECRET=my-secret-key-min-32-characters-long-omgggg
```

### 3. Setup Database

Pastikan Anda sudah membuat database kosong bernama uas (BOLEH SESUAIKAN NAMA DATABASE ABANGDA) di PostgreSQL.

## A. Jalankan Migration (Struktur Tabel)

Gunakan perintah berikut untuk membuat tabel-tabel database:
migrate -path database/migrations -database "postgres://postgres:root@localhost:5432/uas?sslmode=disable" up

## B. Jalankan Seeder (Data Dummy):

Import data awal dari file seedernya.sql yang ada di folder database.

- Buka file database/seeder/seed.sql
- Execute/run query satu per satu dari paling atas.
- Server akan berjalan di http://localhost:3000.
