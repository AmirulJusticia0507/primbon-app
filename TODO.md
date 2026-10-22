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
- [x] Tandai halaman dan navigasi “Hari Baik” dengan status “Segera hadir” karena kalkulasi belum tersedia.
- [x] Tampilkan neptu masing-masing, total, dan sisa perhitungan pada halaman hasil; pastikan pemetaan sisa 0–6 terdokumentasi dan benar menurut aturan yang dipilih.
- [x] Tambahkan navigasi dasar, label form yang terhubung ke input, dan pesan validasi/error yang mudah dipahami.
- [x] Putuskan apakah static assets memang diperlukan: mapping `/static` yang belum digunakan dihapus; UI tetap menggunakan CDN.
- [x] Tambahkan logging dan shutdown/penutupan koneksi database dengan tertib.

## Operasional dan kualitas

- [x] Tambahkan konfigurasi Fiber yang eksplisit, endpoint health check, serta penanganan 404/500.
- [x] Tambahkan pipeline CI untuk `gofmt`, `go vet ./...`, dan `go test ./...`.
- [x] Dokumentasikan kebutuhan konektivitas CDN untuk dependensi UI; AdminLTE tetap disajikan dari CDN.
- [x] Lengkapi dokumentasi skema, langkah pengembangan, dan troubleshooting untuk skema/seed.

## Temuan saat peninjauan

- `go test ./...` berhasil pada saat peninjauan, tetapi tidak ada test di paket mana pun.
- Konfigurasi koneksi PostgreSQL sekarang menggunakan environment variable; startup melaporkan error jika konfigurasi atau koneksi gagal.
- Skema dan seed `petung` tersedia di `db/schema.sql` dan `db/seed.sql`.
- Query perhitungan memeriksa error `Scan` dan mengembalikan respons sesuai untuk kombinasi tidak ditemukan atau kegagalan database; halaman hari baik masih placeholder.
