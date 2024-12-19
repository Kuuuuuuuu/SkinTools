package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

type Config struct {
	SkinPath   string `json:"skinPath"`
	OutputFile string `json:"outputFile"`
	Name       string `json:"name"`
	Geometry   string `json:"geometry"`
	SkinType   string `json:"skinType"`
	NameLength int    `json:"nameLength"`
}

type Skin struct {
	LocalizationName string `json:"localization_name"`
	Geometry         string `json:"geometry"`
	Texture          string `json:"texture"`
	Type             string `json:"type"`
}

type OutputJSON struct {
	Skins            []Skin `json:"skins"`
	SerializeName    string `json:"serialize_name"`
	LocalizationName string `json:"localization_name"`
}

func loadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func generateName(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		result[i] = chars[num.Int64()]
	}
	return string(result)
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	files, err := os.ReadDir(config.SkinPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var skins []Skin
	count := 1

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			oldPath := filepath.Join(config.SkinPath, file.Name())
			newName := generateName(config.NameLength) + ".png"
			newPath := filepath.Join(config.SkinPath, newName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Println("Error renaming file:", err)
				continue
			}

			skin := Skin{
				LocalizationName: fmt.Sprintf("Skin %d", count),
				Geometry:         config.Geometry,
				Texture:          newName,
				Type:             config.SkinType,
			}

			skins = append(skins, skin)
			count++
		}
	}

	output := OutputJSON{
		Skins:            skins,
		SerializeName:    config.Name,
		LocalizationName: config.Name,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile(config.OutputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}

	fmt.Println("Finished with", count-1, "skins")
}
