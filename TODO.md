# TODO / Analisis Kekurangan

Daftar ini mencerminkan kondisi repositori saat ditinjau. Urutan menunjukkan prioritas praktis.

## Prioritas tinggi

- [x] Pindahkan konfigurasi PostgreSQL (host, port, user, password, nama database, SSL) ke environment variable; validasi konfigurasi saat startup dan hindari kredensial default di source.
- [x] Tambahkan skema database dan seed yang terdokumentasi/berulang untuk tabel `petung`, termasuk 35 kombinasi hari-pasaran dan validasi nilai neptu.
- [x] Periksa semua error `QueryRow(...).Scan(...)` dalam handler. Kembalikan respons kesalahan yang sesuai saat database gagal atau kombinasi tidak tersedia.
- [x] Validasi input hari dan pasaran menggunakan daftar nilai yang diizinkan; jangan menghitung hasil dari input kosong/tidak valid.
- [x] Tambahkan unit test pemetaan sisa neptu serta test handler untuk input valid, input tidak valid, dan kondisi gagal database.

## Prioritas menengah

- [x] Pisahkan kalkulasi domain dari handler agar logika mudah diuji; gunakan model `Petung`/`HasilPrimbon` atau hapus tipe yang tidak dipakai.
- [x] Tambahkan kalkulasi weton dan neptu tanggal pada halaman “Hari Baik”; tidak memberi rekomendasi tanpa kaidah tradisi yang dipilih.
- [x] Tambahkan pesan reflektif tentang ikhtiar, tanggung jawab, dan mengambil hikmah; weton ditegaskan bukan penentu takdir.
- [x] Tampilkan neptu masing-masing, total, dan sisa perhitungan pada halaman hasil; pastikan pemetaan sisa 0–6 terdokumentasi dan benar menurut aturan yang dipilih.
- [x] Tambahkan navigasi dasar, label form yang terhubung ke input, dan pesan validasi/error yang mudah dipahami.
- [x] Putuskan apakah static assets memang diperlukan: static lokal digunakan untuk favicon; mapping `/static` menyajikan aset aplikasi.
- [x] Tambahkan logging dan shutdown/penutupan koneksi database dengan tertib.

## Operasional dan kualitas

- [x] Tambahkan konfigurasi Fiber yang eksplisit, endpoint health check, serta penanganan 404/500.
- [x] Tambahkan pipeline CI untuk `gofmt`, `go vet ./...`, dan `go test ./...`.
- [x] Dokumentasikan kebutuhan konektivitas CDN untuk dependensi UI; AdminLTE tetap disajikan dari CDN.
- [x] Lengkapi dokumentasi skema, langkah pengembangan, dan troubleshooting untuk skema/seed.
- [ ] Integrasi Griphub chat telah dibuat sebagai fitur opsional; validasi dengan kredensial/model nyata setelah operator menyediakan keduanya.
- [x] Sediakan `.env` lokal yang diabaikan Git dan `.env.example`; muat konfigurasi dotenv saat startup.
- [x] Sediakan perintah tunggal `task` (default) yang menjalankan pemeriksaan lalu build, di samping `task dev`, `task build`, dan `task check`.

## Temuan saat peninjauan

- `go test ./...` berhasil pada saat peninjauan, tetapi tidak ada test di paket mana pun.
- Konfigurasi koneksi PostgreSQL sekarang menggunakan environment variable; startup melaporkan error jika konfigurasi atau koneksi gagal.
- Skema dan seed `petung` tersedia di `db/schema.sql` dan `db/seed.sql`.
- Query perhitungan memeriksa error `Scan` dan mengembalikan respons sesuai untuk kombinasi tidak ditemukan atau kegagalan database; halaman hari baik masih placeholder.
