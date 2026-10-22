# Primbon App

Aplikasi web sederhana berbasis Go untuk menghitung kecocokan pasangan menurut neptu weton (hari dan pasaran) menggunakan data PostgreSQL.

> Hasil perhitungan adalah informasi budaya/tradisi dan bukan nasihat atau prediksi ilmiah.

## Status fitur

- Beranda: formulir hari dan pasaran dua orang.
- Ramalan jodoh: menghitung jumlah neptu dan memetakan sisa pembagian dengan 7 ke kategori hasil.
- Hari baik: halaman berstatus “Segera hadir”; kalkulasi dan rekomendasi hari belum tersedia.

## Teknologi

- Go 1.24.5 (sesuai `go.mod`)
- Fiber v2
- PostgreSQL melalui `database/sql` dan `lib/pq`
- Template HTML Go; AdminLTE dimuat dari CDN

## Prasyarat

1. Go 1.24.5 atau versi kompatibel.
2. PostgreSQL berjalan dan dapat diakses aplikasi.
3. Database memiliki skema dan data `petung` yang disiapkan dari `db/schema.sql` dan `db/seed.sql`.

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

Siapkan tabel dan data awal secara berurutan (seed aman dijalankan ulang):

```sh
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -f db/schema.sql
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -f db/seed.sql
```

Di PowerShell, jalankan `psql` dengan parameter nilai yang sama dari environment bila ekspansi variabel gaya shell POSIX tidak tersedia.

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

- Pemetaan hasil memakai `total neptu % 7`: sisa 0 Lebu Ketiup Angin, 1 Wasesa Segara, 2 Tunggak Semi, 3 Satria Wibawa, 4 Sumur Sinaba, 5 Satria Wirang, dan 6 Bumi Kepetak. Ini mengikuti konvensi yang digunakan aplikasi dan merupakan tradisi budaya, bukan prediksi ilmiah.

- Handler perhitungan saat ini mengabaikan error hasil query. Bila kombinasi tidak ditemukan, neptu dapat menjadi nol dan hasil tetap ditampilkan.
- `models/primbon.go` mendefinisikan struktur, tetapi belum digunakan oleh handler.
- Belum ada kebutuhan aset lokal khusus; mapping `/static` yang tidak digunakan telah dihapus. AdminLTE tetap dimuat dari CDN.
- Hasil ramalan hanya menampilkan total neptu; nilai individual dan sisa belum ditampilkan.
- Baca [`TODO.md`](TODO.md) untuk kekurangan dan prioritas perbaikan. Panduan kontribusi/agen ada di [`AGENTS.md`](AGENTS.md).
