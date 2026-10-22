package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHitungHasilSemuaSisa(t *testing.T) {
	cases := []struct {
		sisa  int
		hasil string
	}{
		{0, "Lebu Ketiup Angin"},
		{1, "Wasesa Segara"},
		{2, "Tunggak Semi"},
		{3, "Satria Wibawa"},
		{4, "Sumur Sinaba"},
		{5, "Satria Wirang"},
		{6, "Bumi Kepetak"},
	}
	for _, tc := range cases {
		t.Run(tc.hasil, func(t *testing.T) {
			total, sisa, hasil := HitungHasil(tc.sisa, 0)
			if total != tc.sisa || sisa != tc.sisa || !strings.HasPrefix(hasil, tc.hasil) {
				t.Fatalf("HitungHasil(%d, 0) = (%d, %d, %q)", tc.sisa, total, sisa, hasil)
			}
		})
	}
}

func TestRamalJodohInputTidakValid(t *testing.T) {
	app := fiber.New()
	app.Post("/ramal", func(c *fiber.Ctx) error { return RamalJodoh(c, nil) })
	request := httptest.NewRequest("POST", "/ramal", strings.NewReader("hari_laki=Invalid"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusBadRequest)
	}
}
