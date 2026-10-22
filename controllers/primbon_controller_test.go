package controllers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

type fakePetungDriver struct{}
type fakePetungConn struct{ state *fakePetungState }
type fakePetungRows struct {
	value int
	done  bool
}

type fakePetungState struct {
	mu      sync.Mutex
	values  []int
	errors  []error
	queries int
}

var fakePetungStates sync.Map

func init() {
	sql.Register("fake-petung", fakePetungDriver{})
}

func (fakePetungDriver) Open(name string) (driver.Conn, error) {
	state, ok := fakePetungStates.Load(name)
	if !ok {
		return nil, fmt.Errorf("state database palsu %q tidak ditemukan", name)
	}
	return &fakePetungConn{state: state.(*fakePetungState)}, nil
}

func (c *fakePetungConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("Prepare tidak didukung")
}
func (c *fakePetungConn) Close() error { return nil }
func (c *fakePetungConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transaksi tidak didukung")
}
func (c *fakePetungConn) QueryContext(_ context.Context, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	c.state.mu.Lock()
	defer c.state.mu.Unlock()
	index := c.state.queries
	c.state.queries++
	if index < len(c.state.errors) && c.state.errors[index] != nil {
		return nil, c.state.errors[index]
	}
	if index >= len(c.state.values) {
		return nil, errors.New("tidak ada hasil palsu tersedia")
	}
	return &fakePetungRows{value: c.state.values[index]}, nil
}

func (*fakePetungRows) Columns() []string { return []string{"neptu"} }
func (*fakePetungRows) Close() error      { return nil }
func (r *fakePetungRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	dest[0] = int64(r.value)
	r.done = true
	return nil
}

func openFakePetungDB(t *testing.T, state *fakePetungState) *sql.DB {
	t.Helper()
	name := fmt.Sprintf("test-%p", state)
	fakePetungStates.Store(name, state)
	t.Cleanup(func() { fakePetungStates.Delete(name) })
	db, err := sql.Open("fake-petung", name)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

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

func TestWetonTanggalAnchorAndFollowingDay(t *testing.T) {
	friday := time.Date(2024, time.January, 5, 0, 0, 0, 0, time.UTC)
	day, market, dayValue, marketValue := WetonTanggal(friday)
	if day != "Jumat" || market != "Legi" || dayValue != 6 || marketValue != 5 {
		t.Fatalf("anchor = %s %s (%d,%d), want Jumat Legi (6,5)", day, market, dayValue, marketValue)
	}
	saturday := friday.AddDate(0, 0, 1)
	day, market, _, _ = WetonTanggal(saturday)
	if day != "Sabtu" || market != "Pahing" {
		t.Fatalf("following day = %s %s, want Sabtu Pahing", day, market)
	}
}

func TestRamalJodohInputTidakValid(t *testing.T) {
	withProjectRoot(t)
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

func TestRamalJodohInputValid(t *testing.T) {
	withProjectRoot(t)
	db := openFakePetungDB(t, &fakePetungState{values: []int{9, 9}})
	app := fiber.New()
	app.Post("/ramal", func(c *fiber.Ctx) error { return RamalJodoh(c, db) })
	response := postWeton(t, app)
	if response.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusOK)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Orang pertama", "Orang kedua", "Total neptu", "18", "Sisa dibagi 7", "Sumur Sinaba"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("response body tidak memuat %q", want)
		}
	}
}

func TestRamalJodohDatabaseError(t *testing.T) {
	withProjectRoot(t)
	db := openFakePetungDB(t, &fakePetungState{errors: []error{errors.New("database tidak tersedia")}})
	app := fiber.New()
	app.Post("/ramal", func(c *fiber.Ctx) error { return RamalJodoh(c, db) })
	response := postWeton(t, app)
	if response.StatusCode != fiber.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusInternalServerError)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "Data perhitungan tidak dapat diakses") {
		t.Errorf("pesan respons tidak menjelaskan kegagalan database: %s", body)
	}
	if strings.Contains(string(body), "database tidak tersedia") {
		t.Error("detail error database tidak boleh dikirim ke pengguna")
	}
}

func postWeton(t *testing.T, app *fiber.App) *http.Response {
	t.Helper()
	request := httptest.NewRequest("POST", "/ramal", strings.NewReader("hari_laki=Senin&pasaran_laki=Legi&hari_perempuan=Senin&pasaran_perempuan=Legi"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func withProjectRoot(t *testing.T) {
	t.Helper()
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(".."); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(workingDir) })
}
