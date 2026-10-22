package controllers

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"primbon-app/internal/assets"
)

func renderHTML(c *fiber.Ctx, tmpl string, data interface{}) error {
	files, err := fs.Sub(assets.Files, "views")
	if err != nil {
		return err
	}
	t, err := template.ParseFS(files, "layouts/base.html", strings.TrimPrefix(tmpl, "views/"))
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
		c.Status(fiber.StatusBadRequest)
		return renderHTML(c, "views/home.html", fiber.Map{
			"Title": "Ramalan Primbon Jawa",
			"Error": "Pilih hari dan pasaran yang valid untuk kedua orang.",
		})
	}

	var neptuL, neptuP int
	if err := db.QueryRow(`SELECT nilai_hari + nilai_pasaran FROM petung WHERE hari=$1 AND pasaran=$2`, hariL, pasaranL).Scan(&neptuL); err != nil {
		return petungQueryError(c, err)
	}
	if err := db.QueryRow(`SELECT nilai_hari + nilai_pasaran FROM petung WHERE hari=$1 AND pasaran=$2`, hariP, pasaranP).Scan(&neptuP); err != nil {
		return petungQueryError(c, err)
	}

	total, sisa, hasil := HitungHasil(neptuL, neptuP)
	return renderHTML(c, "views/ramalan.html", fiber.Map{
		"Title":  "Hasil Ramalan Jodoh",
		"Hasil":  hasil,
		"NeptuL": neptuL,
		"NeptuP": neptuP,
		"Total":  total,
		"Sisa":   sisa,
	})
}

func HitungHasil(neptuL, neptuP int) (total, sisa int, hasil string) {
	total = neptuL + neptuP
	sisa = total % 7
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

	return total, sisa, hasil
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
		c.Status(fiber.StatusBadRequest)
		return renderHTML(c, "views/home.html", fiber.Map{
			"Title": "Ramalan Primbon Jawa",
			"Error": "Kombinasi hari dan pasaran tidak ditemukan. Pastikan data petung sudah disiapkan.",
		})
	}
	log.Printf("gagal mengambil data petung: %v", err)
	c.Status(fiber.StatusInternalServerError)
	return renderHTML(c, "views/home.html", fiber.Map{
		"Title": "Ramalan Primbon Jawa",
		"Error": "Data perhitungan tidak dapat diakses saat ini. Silakan coba lagi nanti.",
	})
}

func HariBaik(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Hitung Hari Baik",
		"Today": time.Now().Format("2006-01-02"),
	}
	if c.Method() == fiber.MethodPost {
		date, err := time.Parse("2006-01-02", c.FormValue("tanggal"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			data["Error"] = "Masukkan tanggal yang valid."
			return renderHTML(c, "views/hari_baik.html", data)
		}
		hari, pasaran, neptuHari, neptuPasaran := WetonTanggal(date)
		data["Tanggal"] = date.Format("Monday, 2 January 2006")
		data["SelectedDate"] = date.Format("2006-01-02")
		data["Hari"] = hari
		data["Pasaran"] = pasaran
		data["NeptuHari"] = neptuHari
		data["NeptuPasaran"] = neptuPasaran
		data["Total"] = neptuHari + neptuPasaran
	}
	return renderHTML(c, "views/hari_baik.html", data)
}

// WetonTanggal computes the Javanese weekday and five-day market cycle.
// The pasaran anchor is Friday Legi on 2024-01-05.
func WetonTanggal(date time.Time) (hari, pasaran string, neptuHari, neptuPasaran int) {
	hariNames := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	hariNeptu := []int{5, 4, 3, 7, 8, 6, 9}
	pasaranNames := []string{"Legi", "Pahing", "Pon", "Wage", "Kliwon"}
	pasaranNeptu := []int{5, 9, 7, 4, 8}

	day := date.UTC().Weekday()
	dayIndex := int(day)
	anchor := time.Date(2024, time.January, 5, 0, 0, 0, 0, time.UTC)
	selected := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	daysFromAnchor := int(selected.Sub(anchor).Hours() / 24)
	marketIndex := ((daysFromAnchor % 5) + 5) % 5
	return hariNames[dayIndex], pasaranNames[marketIndex], hariNeptu[dayIndex], pasaranNeptu[marketIndex]
}
