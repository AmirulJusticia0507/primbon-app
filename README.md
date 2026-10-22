# Primbon App

Aplikasi web sederhana berbasis Go untuk menghitung kecocokan pasangan menurut neptu weton (hari dan pasaran) menggunakan data PostgreSQL.

> Hasil perhitungan adalah informasi budaya/tradisi dan bukan nasihat atau prediksi ilmiah.

## Status fitur

- Beranda: formulir hari dan pasaran dua orang.
- Ramalan jodoh: menghitung jumlah neptu dan memetakan sisa pembagian dengan 7 ke kategori hasil.
- Hari baik: halaman placeholder; kalkulasi belum tersedia.

## Teknologi

- Go 1.24.5 (sesuai `go.mod`)
- Fiber v2
- PostgreSQL melalui `database/sql` dan `lib/pq`
- Template HTML Go; AdminLTE dimuat dari CDN

## Prasyarat

1. Go 1.24.5 atau versi kompatibel.
2. PostgreSQL berjalan dan dapat diakses aplikasi.
3. Database memiliki tabel `petung` berisi kolom `hari`, `pasaran`, `nilai_hari`, dan `nilai_pasaran`, dengan satu baris untuk setiap kombinasi hari/pasaran (35 kombinasi). Belum tersedia berkas migrasi atau seed; siapkan skema dan data tersebut sebelum menjalankan aplikasi.

## Menjalankan

Konfigurasi PostgreSQL wajib diberikan melalui environment variable berikut (tidak ada kredensial default di source):

| Variabel | Keterangan | Contoh |
| --- | --- | --- |
| `DB_HOST` | Host PostgreSQL | `localhost` |
| `DB_PORT` | Port (1–65535) | `5432` |
| `DB_USER` | Pengguna database | `postgres` |
| `DB_PASSWORD` | Kata sandi database | *(isi lokal, jangan commit)* |
| `DB_NAME` | Nama database | `primbon_db` |
| `DB_SSLMODE` | Mode SSL PostgreSQL (`disable`, `require`, `verify-ca`, `verify-full`, `allow`, `prefer`) | `disable` |

Semua variabel divalidasi saat startup; aplikasi berhenti dengan pesan kesalahan jika konfigurasi tidak lengkap/tidak valid atau database tidak dapat dihubungi.

Contoh PowerShell untuk sesi terminal saat ini:

```powershell
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "isi-password-lokal"
$env:DB_NAME = "primbon_db"
$env:DB_SSLMODE = "disable"
```

```sh
go mod download
go run .
```

Server mendengarkan pada `http://localhost:8080`. Buka `/` untuk formulir dan `/hari-baik` untuk halaman fitur yang belum selesai. Pastikan PostgreSQL dan tabel `petung` siap lebih dahulu; aplikasi akan berhenti jika koneksi database gagal.

## Pemeriksaan

```sh
go test ./...
```

Saat dokumentasi ini dibuat, perintah berhasil tetapi seluruh paket belum memiliki file test.

## Struktur proyek

```text
.
├── main.go                         # Inisialisasi database dan server
├── config/database.go              # Konfigurasi koneksi PostgreSQL
├── routes/routes.go                # Rute Fiber
├── controllers/primbon_controller.go # Handler halaman dan perhitungan
├── models/primbon.go                # Struktur data
└── views/                            # Template HTML
```

## Catatan implementasi

- Handler perhitungan saat ini mengabaikan error hasil query. Bila kombinasi tidak ditemukan, neptu dapat menjadi nol dan hasil tetap ditampilkan.
- `models/primbon.go` mendefinisikan struktur, tetapi belum digunakan oleh handler.
- Static file diarahkan ke `./static`, namun direktori tersebut belum ada dalam repositori.
- Hasil ramalan hanya menampilkan total neptu; nilai individual dan sisa belum ditampilkan.
- Baca [`TODO.md`](TODO.md) untuk kekurangan dan prioritas perbaikan. Panduan kontribusi/agen ada di [`AGENTS.md`](AGENTS.md).
