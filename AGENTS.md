# Panduan Agen dan Kontributor

Panduan ini berlaku untuk seluruh proyek.

## Konteks proyek

- Aplikasi Go dengan Fiber v2 dan PostgreSQL (`database/sql`, driver `lib/pq`).
- Template HTML berada di `views/`; entry point ada di `main.go`.
- Bahasa antarmuka dan dokumentasi pengguna adalah Bahasa Indonesia.
- Alur utama: `routes` memanggil `controllers`, handler menghitung neptu dari tabel PostgreSQL `petung` dan merender template.

## Aturan perubahan

1. Ikuti pola dan gaya Go yang sudah ada; jalankan `gofmt` pada berkas Go yang berubah.
2. Gunakan query berparameter, jangan menyusun SQL dari input pengguna.
3. Jangan menambahkan kredensial atau rahasia ke source control. Konfigurasi database sebaiknya melalui environment variable saat refactor.
4. Tangani error query, input, dan rendering secara eksplisit; jangan diam-diam mengubah kegagalan query menjadi hasil valid.
5. Pertahankan bahasa Indonesia untuk teks UI. Jelaskan secara proporsional bahwa hasil primbon merupakan tradisi budaya.
6. Jangan mengklaim fungsi selesai jika masih placeholder. Perbarui `README.md` dan `TODO.md` ketika status fitur atau cara menjalankan berubah.
7. Perubahan skema database perlu disertai migrasi dan/atau instruksi seed yang dapat diulang.

## Verifikasi

Untuk perubahan Go, jalankan:

```sh
gofmt -w <berkas-go-yang-diubah>
go test ./...
```

Jika perubahan menyentuh perilaku database, verifikasi juga dengan PostgreSQL dan data `petung` yang valid. Jangan menyatakan pengujian integrasi berhasil jika database tersebut tidak tersedia.

## Struktur penting

- `config/`: konfigurasi dan inisialisasi koneksi.
- `routes/`: deklarasi rute HTTP.
- `controllers/`: handler dan logika aplikasi.
- `models/`: tipe data domain.
- `views/`: template HTML.
- `TODO.md`: daftar pekerjaan yang diketahui.
