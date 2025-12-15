# Sistem Pelaporan Prestasi Mahasiswa (Backend API)

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![Fiber Framework](https://img.shields.io/badge/Fiber-v2-black?style=flat&logo=gofiber)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)
![MongoDB](https://img.shields.io/badge/MongoDB-6.0+-47A248?style=flat&logo=mongodb)

Backend API untuk **Sistem Pelaporan dan Validasi Prestasi Mahasiswa**. Project ini dibangun dengan arsitektur **Hybrid Database**, menggunakan **PostgreSQL** untuk data relasional & kontrol akses (RBAC), serta **MongoDB** untuk penyimpanan data prestasi yang bersifat dinamis.

---

## 📋 Fitur Utama

- **Autentikasi JWT**

  - Login & Refresh Token

- **Role-Based Access Control (RBAC)**

  - Admin
  - Mahasiswa
  - Dosen Wali

- **Manajemen Prestasi Mahasiswa**

  - **Hybrid Storage**

    - PostgreSQL: data referensi, relasi, dan status
    - MongoDB: detail prestasi dinamis

  - **Workflow Status**

    - `draft` → `submitted` → `verified` / `rejected`

  - **Validasi Hak Akses**

    - Dosen Wali hanya dapat memvalidasi mahasiswa bimbingannya

- **Manajemen User & Data Mahasiswa**

---

## 🛠️ Tech Stack

- **Language:** Golang
- **Framework:** Fiber v2
- **Database:**

  - PostgreSQL (Relational Database)
  - MongoDB (NoSQL Database)

- **Migration Tool:** Golang Migrate

---

## 🚀 Cara Menjalankan Project

### 1. Prasyarat

Pastikan tools berikut sudah terinstall:

- [Go](https://go.dev/dl/) (versi 1.20 atau lebih baru)
- [PostgreSQL](https://www.postgresql.org/)
- [MongoDB](https://www.mongodb.com/)
- [Golang Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

---

### 2. Konfigurasi Environment

Buat file **`.env`** di root project (file ini **tidak boleh di-push ke GitHub**) dan isi dengan konfigurasi berikut:

```env
APP_PORT=3000
API_KEY=your-api-key
POSTGRES_URI=postgres://<user>:<password>@<host>:<port>/<db>?sslmode=disable
MONGO_URI=mongodb://<host>:<port>
MONGO_DB=uas
JWT_SECRET=your-secret-key-min-32-characters
```

📌 **Catatan:**

- Pastikan `JWT_SECRET` memiliki panjang minimal 32 karakter.
- Untuk production, gunakan credential yang lebih aman.

---

### 3. Setup Database PostgreSQL

Pastikan Anda sudah membuat database kosong dengan nama **`uas`** (atau sesuaikan dengan nilai `POSTGRES_URI`).

```sql
CREATE DATABASE uas;
```

---

## A. Menjalankan Migration (Struktur Tabel)

Gunakan perintah berikut untuk menjalankan migration database:

```bash
migrate -path database/migrations \
  -database "postgres://postgres:root@localhost:5432/uas?sslmode=disable" up
```

Migration ini akan membuat:

- tabel users
- tabel students
- tabel achievement_references
- enum status prestasi
- relasi antar tabel

---

## B. Menjalankan Seeder (Data Awal)

Untuk mengisi data awal (dummy data):

1. Buka file berikut:

   ```
   database/seeder/seed.sql
   ```

2. Jalankan query **secara berurutan dari atas ke bawah** menggunakan:

   - pgAdmin, atau
   - psql CLI, atau
   - database client lainnya

---

## ▶️ Menjalankan Server

Setelah migration dan seeder selesai, jalankan aplikasi:

```bash
go run main.go
```

Server akan berjalan di:

```
http://localhost:3000
```

---

## 📌 Catatan Tambahan

- Project ini menggunakan **arsitektur repository pattern**.
- MongoDB akan otomatis membuat collection saat data pertama kali di-insert.
- File `.env` wajib dimasukkan ke `.gitignore`.
- Gunakan `.env.example` sebagai template konfigurasi.

---

## 📄 Lisensi

Project ini dibuat untuk keperluan pembelajaran dan pengembangan internal.
