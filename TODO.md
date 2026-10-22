# TODO / Analisis Kekurangan

Daftar ini mencerminkan kondisi repositori saat ditinjau. Urutan menunjukkan prioritas praktis.

## Prioritas tinggi

- [x] Pindahkan konfigurasi PostgreSQL (host, port, user, password, nama database, SSL) ke environment variable; validasi konfigurasi saat startup dan hindari kredensial default di source.
- [x] Tambahkan skema database dan seed yang terdokumentasi/berulang untuk tabel `petung`, termasuk 35 kombinasi hari-pasaran dan validasi nilai neptu.
- [x] Periksa semua error `QueryRow(...).Scan(...)` dalam handler. Kembalikan respons kesalahan yang sesuai saat database gagal atau kombinasi tidak tersedia.
- [ ] Validasi input hari dan pasaran menggunakan daftar nilai yang diizinkan; jangan menghitung hasil dari input kosong/tidak valid.
- [ ] Tambahkan unit test untuk pemetaan sisa neptu ke hasil serta test handler untuk input valid dan kondisi gagal database.

## Prioritas menengah

- [ ] Pisahkan kalkulasi domain dari handler agar logika mudah diuji; gunakan model `Petung`/`HasilPrimbon` atau hapus tipe yang tidak dipakai.
- [ ] Tambahkan halaman “Hari Baik” yang benar-benar menghitung hasil atau ubah penamaan navigasinya agar status coming soon jelas.
- [ ] Tampilkan neptu masing-masing, total, dan sisa perhitungan pada halaman hasil; pastikan pemetaan sisa 0–6 terdokumentasi dan benar menurut aturan yang dipilih.
- [ ] Tambahkan navigasi dasar, label form yang terhubung ke input, dan pesan validasi/error yang mudah dipahami.
- [ ] Putuskan apakah static assets memang diperlukan: tambahkan direktori `static/` dan aset lokal, atau hapus mapping `/static` yang belum digunakan.
- [ ] Tambahkan logging dan shutdown/penutupan koneksi database dengan tertib.

## Operasional dan kualitas

- [ ] Tambahkan konfigurasi Fiber yang eksplisit, endpoint health check, serta penanganan 404/500.
- [ ] Tambahkan pemeriksaan dependency dan pipeline CI untuk `gofmt`, `go vet ./...`, dan `go test ./...`.
- [ ] Pertimbangkan menyimpan dependensi UI secara lokal atau mendokumentasikan kebutuhan konektivitas CDN.
- [ ] Lengkapi dokumentasi skema, langkah pengembangan, dan troubleshooting setelah migrasi/seed tersedia.

## Temuan saat peninjauan

- `go test ./...` berhasil pada saat peninjauan, tetapi tidak ada test di paket mana pun.
- Konfigurasi koneksi PostgreSQL sekarang menggunakan environment variable; startup melaporkan error jika konfigurasi atau koneksi gagal.
- Skema dan seed `petung` tersedia di `db/schema.sql` dan `db/seed.sql`.
- Query perhitungan memeriksa error `Scan` dan mengembalikan respons sesuai untuk kombinasi tidak ditemukan atau kegagalan database; halaman hari baik masih placeholder.
