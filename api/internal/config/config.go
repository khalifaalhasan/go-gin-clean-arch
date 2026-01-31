package config

import (
	"github.com/spf13/viper"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	GinMode       string `mapstructure:"GIN_MODE"`
}

// LoadConfig membaca konfigurasi dari file ATAU environment variable
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)  // Cari file config di folder path
	viper.SetConfigName("app") // Nama file: app.env
	viper.SetConfigType("env") // Tipe file: env

	viper.AutomaticEnv() // BACA ENV VARIABLE (Penting buat Render!)

	err = viper.ReadInConfig()
	if err != nil {
		// Kalau file app.env tidak ditemukan (kasus di Render),
		// tidak apa-apa asalkan env variable sudah diset.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}