package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

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

const (
	skinPath   = "./images"
	outputFile = "skins.json"
	name       = "IDK"
)

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
	files, err := os.ReadDir(skinPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var skins []Skin
	count := 1

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".png" {
			oldPath := filepath.Join(skinPath, file.Name())
			newName := generateName(16) + ".png"
			newPath := filepath.Join(skinPath, newName)

			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Println("Error renaming file:", err)
				continue
			}

			skin := Skin{
				LocalizationName: fmt.Sprintf("Skin %d", count),
				Geometry:         "geometry.humanoid.customSlim", // change this if u want
				Texture:          newName,
				Type:             "free",
			}

			skins = append(skins, skin)
			count++
		}
	}

	output := OutputJSON{
		Skins:            skins,
		SerializeName:    name,
		LocalizationName: name,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}

	fmt.Println("Finished with", count-1, "skins")
}
