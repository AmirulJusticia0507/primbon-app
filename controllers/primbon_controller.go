package controllers

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"
)

func renderHTML(c *fiber.Ctx, tmpl string, data interface{}) error {
	t, err := template.ParseFiles("views/layouts/base.html", tmpl)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "base", data); err != nil {
		return err
	}
	return c.Type("html").SendString(buf.String())
}

func Home(c *fiber.Ctx) error {
	return renderHTML(c, "views/home.html", fiber.Map{
		"Title": "Ramalan Primbon Jawa",
	})
}

func RamalJodoh(c *fiber.Ctx, db *sql.DB) error {
	hariL := c.FormValue("hari_laki")
	pasaranL := c.FormValue("pasaran_laki")
	hariP := c.FormValue("hari_perempuan")
	pasaranP := c.FormValue("pasaran_perempuan")
	if !hariValid(hariL) || !hariValid(hariP) || !pasaranValid(pasaranL) || !pasaranValid(pasaranP) {
		return c.Status(fiber.StatusBadRequest).SendString("Pilih hari dan pasaran yang valid untuk kedua orang.")
	}

	var neptuL, neptuP int
	if err := db.QueryRow(`SELECT nilai_hari + nilai_pasaran FROM petung WHERE hari=$1 AND pasaran=$2`, hariL, pasaranL).Scan(&neptuL); err != nil {
		return petungQueryError(c, err)
	}
	if err := db.QueryRow(`SELECT nilai_hari + nilai_pasaran FROM petung WHERE hari=$1 AND pasaran=$2`, hariP, pasaranP).Scan(&neptuP); err != nil {
		return petungQueryError(c, err)
	}

	total := neptuL + neptuP
	sisa := total % 7
	var hasil string

	switch sisa {
	case 1:
		hasil = "Wasesa Segara (baik, pasangan saling melengkapi)"
	case 2:
		hasil = "Tunggak Semi (baik, rezeki berkembang)"
	case 3:
		hasil = "Satria Wibawa (baik, berwibawa dan dihormati)"
	case 4:
		hasil = "Sumur Sinaba (baik, rumah tangga sejahtera)"
	case 5:
		hasil = "Satria Wirang (kurang baik, sering ada salah paham)"
	case 6:
		hasil = "Bumi Kepetak (kurang harmonis)"
	default:
		hasil = "Lebu Ketiup Angin (buruk, sering bertengkar)"
	}

	return renderHTML(c, "views/ramalan.html", fiber.Map{
		"Title":  "Hasil Ramalan Jodoh",
		"Hasil":  hasil,
		"NeptuL": neptuL,
		"NeptuP": neptuP,
		"Total":  total,
		"Sisa":   sisa,
	})
}

func hariValid(hari string) bool {
	switch hari {
	case "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu":
		return true
	default:
		return false
	}
}

func pasaranValid(pasaran string) bool {
	switch pasaran {
	case "Legi", "Pahing", "Pon", "Wage", "Kliwon":
		return true
	default:
		return false
	}
}

func petungQueryError(c *fiber.Ctx, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusBadRequest).SendString("Kombinasi hari dan pasaran tidak ditemukan. Periksa kembali data weton.")
	}
	log.Printf("gagal mengambil data petung: %v", err)
	return c.Status(fiber.StatusInternalServerError).SendString("Terjadi kesalahan saat menghitung weton. Silakan coba lagi nanti.")
}

func HariBaik(c *fiber.Ctx) error {
	return renderHTML(c, "views/hari_baik.html", fiber.Map{
		"Title": "Hitung Hari Baik",
	})
}
