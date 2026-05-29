package utils

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"regexp"
)

// Map untuk menyimpan data di RAM supaya pencarian super cepat (O(1))
var cityCodeMap map[string]string

// Struct disesuaikan persis dengan skema JSON baru
type Regency struct {
	ID         string  `json:"id"`
	ProvinceID string  `json:"province_id"`
	Name       string  `json:"name"`
	AltName    string  `json:"alt_name"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

// Fungsi init akan otomatis berjalan sekali saat aplikasi pertama kali start
func init() {
	cityCodeMap = make(map[string]string)

	// Pastikan file regencies.json berada di root folder proyek (sejajar main.go)
	file, err := os.ReadFile("regencies.json")
	if err != nil {
		log.Println("Peringatan: file regencies.json tidak ditemukan, parser NIK dilewati")
		return
	}

	var regencies []Regency
	err = json.Unmarshal(file, &regencies)
	if err != nil {
		log.Println("Peringatan: Gagal membaca format JSON:", err)
		return
	}

	// Memindahkan data dari slice ke map untuk optimasi performa
	for _, reg := range regencies {
		cityCodeMap[reg.ID] = reg.Name
	}
}

// GetCityFromNIK mengambil nama Kota/Kabupaten berdasarkan 4 digit pertama NIK
func GetCityFromNIK(nik string) (string, error) {
	// 1. Validasi panjang karakter
	if len(nik) != 16 {
		return "", errors.New("format tidak valid: panjang NIK harus 16 digit")
	}

	// 2. Validasi isi karakter (harus angka semua)
	isNumeric := regexp.MustCompile(`^[0-9]+$`).MatchString(nik)
	if !isNumeric {
		return "", errors.New("format tidak valid: NIK hanya boleh berisi angka")
	}

	// 3. Potong 4 digit pertama (Kode Provinsi + Kode Kota/Kabupaten)
	cityCode := nik[0:4]

	// 4. Cari data di memory map
	city, exists := cityCodeMap[cityCode]
	if !exists {
		return "Luar Wilayah / Tidak Diketahui", nil
	}

	return city, nil
}
