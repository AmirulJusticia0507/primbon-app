package models

type Petung struct {
	ID           int
	Hari         string
	NilaiHari    int
	Pasaran      string
	NilaiPasaran int
}

type HasilPrimbon struct {
	NamaLaki       string
	WetonLaki      string
	NeptuLaki      int
	NamaPerempuan  string
	WetonPerempuan string
	NeptuPerempuan int
	TotalNeptu     int
	Sisa           int
	Hasil          string
}
