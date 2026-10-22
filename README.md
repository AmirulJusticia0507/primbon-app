# Primbon App

Aplikasi web sederhana berbasis Go untuk menghitung kecocokan pasangan menurut neptu weton (hari dan pasaran) menggunakan data PostgreSQL.

> Hasil perhitungan adalah informasi budaya/tradisi dan bukan nasihat atau prediksi ilmiah.

## Status fitur

- Beranda: formulir hari dan pasaran dua orang.
- Ramalan jodoh: menghitung jumlah neptu dan memetakan sisa pembagian dengan 7 ke kategori hasil.
- Hari baik: halaman berstatus “Segera hadir”; kalkulasi dan rekomendasi hari belum tersedia.
- Tanya AI: fitur chat opsional melalui Griphub Router; perlu API key dan slug model dari operator layanan.

## Teknologi

- Go 1.24.5 (sesuai `go.mod`)
- Fiber v2
- PostgreSQL melalui `database/sql` dan `lib/pq`
- Template HTML Go; AdminLTE dimuat dari CDN jsDelivr

### Konektivitas aset UI

AdminLTE CSS dan Google Fonts dimuat dari CDN jsDelivr/Google Fonts, bukan disimpan lokal. Pengguna memerlukan konektivitas internet agar library UI dan tipografi tersedia; jika CDN tidak dapat dijangkau, halaman tetap dapat digunakan dengan fallback font sistem dan styling inti lokal, tetapi AdminLTE tidak dimuat. Favicon SVG disajikan lokal dari `static/`.

## Prasyarat

1. Go 1.24.5 atau versi kompatibel.
2. PostgreSQL berjalan dan dapat diakses aplikasi.
3. PostgreSQL menyediakan database target. Tabel dan data `petung` dibuat dengan berkas yang ada di `db/` (lihat langkah di bawah).

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

Fitur Tanya AI tidak wajib untuk menjalankan aplikasi. Untuk mengaktifkannya, atur `AI_API_KEY` dan `AI_MODEL` pada environment proses aplikasi. Jangan masukkan key ke source control. Tanpa konfigurasi tersebut, endpoint AI memberi tahu bahwa fitur belum dikonfigurasi.

Contoh PowerShell untuk sesi terminal saat ini:

```powershell
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "isi-password-lokal"
$env:DB_NAME = "primbon_db"
$env:DB_SSLMODE = "disable"
```

### Skema dan seed

`db/schema.sql` membuat tabel `petung` dengan kolom identity `id`, nama `hari` dan `pasaran`, serta nilai neptu masing-masing. CHECK constraint membatasi nama ke tujuh hari dan lima pasaran serta memastikan setiap nama memiliki nilai neptu yang sesuai. Constraint unik mencegah duplikasi pasangan hari-pasaran. `db/seed.sql` memasukkan semua 35 kombinasi; `ON CONFLICT` memperbarui nilai yang sudah ada sehingga aman dijalankan berulang.

Siapkan tabel dan data awal secara berurutan (seed aman dijalankan ulang):

```sh
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -f db/schema.sql
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 -f db/seed.sql
```

Di PowerShell, gunakan parameter langsung, misalnya `psql -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -d $env:DB_NAME -v ON_ERROR_STOP=1 -f db/schema.sql`, lalu ulangi untuk `db/seed.sql`. `psql` akan meminta password jika belum tersedia melalui mekanisme aman lokal.

Verifikasi jumlah data setelah seed:

```sql
SELECT count(*) FROM petung; -- diharapkan 35
```

```sh
go mod download
go run .
```

Server mendengarkan pada `http://localhost:8080`. Buka `/` untuk formulir dan `/hari-baik` untuk melihat status fitur yang segera hadir. Endpoint `GET /health` memeriksa koneksi database dan mengembalikan HTTP 503 jika tidak sehat. Pastikan PostgreSQL dan tabel `petung` siap lebih dahulu; aplikasi akan berhenti jika koneksi database gagal. Saat menerima SIGINT/SIGTERM server menjalankan graceful shutdown dengan batas waktu 10 detik, kemudian koneksi database ditutup.

## Pemeriksaan

```sh
gofmt -l $(git ls-files '*.go')
go vet ./...
go test ./...
```

CI GitHub Actions menjalankan pemeriksaan `gofmt`, `go vet ./...`, dan `go test ./...` pada push serta pull request. Test mencakup kalkulasi dan handler; beberapa paket belum memiliki test tersendiri.

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

## Pengembangan dan troubleshooting

- Pemetaan hasil memakai `total neptu % 7`: sisa 0 Lebu Ketiup Angin, 1 Wasesa Segara, 2 Tunggak Semi, 3 Satria Wibawa, 4 Sumur Sinaba, 5 Satria Wirang, dan 6 Bumi Kepetak. Ini mengikuti konvensi yang digunakan aplikasi dan merupakan tradisi budaya, bukan prediksi ilmiah.
- Jika aplikasi gagal startup dengan pesan variabel wajib, pastikan keenam `DB_*` (host, port, user, password, nama database, SSL mode) sudah diatur di proses yang menjalankan aplikasi.
- Jika koneksi gagal, pastikan PostgreSQL berjalan, host/port dapat dijangkau, database dan hak akses user tersedia, serta `DB_SSLMODE` sesuai konfigurasi server. Aplikasi memeriksa koneksi saat startup.
- Jika health check `GET /health` memberi HTTP 503, koneksi database saat ini gagal; periksa ketersediaan PostgreSQL dan log aplikasi.
- Jika query tidak menemukan kombinasi weton, pastikan kedua berkas SQL sudah dijalankan dan `SELECT count(*) FROM petung` menghasilkan 35. Seed dapat dijalankan kembali untuk memperbaiki/melengkapi baris.
- Jika tampilan tidak bergaya, browser mungkin tidak dapat menjangkau CDN jsDelivr; HTML tetap tersedia tetapi AdminLTE tidak. Lihat bagian Konektivitas aset UI.
- `models/primbon.go` masih mendefinisikan struktur domain yang belum digunakan handler.

Untuk kekurangan dan prioritas perbaikan, baca [`TODO.md`](TODO.md). Panduan kontribusi/agen ada di [`AGENTS.md`](AGENTS.md).
