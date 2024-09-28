package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const resultsFileName string = "__output-results.csv"

// searchInFiles busca un término en todos los archivos de un directorio y guarda las coincidencias en results.txt
func searchInFiles(directory, searchTerm string) error {
	// Abrir o crear el archivo donde se guardarán los resultados
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error leer el directorio: %v", err)
	}

	resultsFilePath := cwd + "/" + resultsFileName

	fmt.Printf("Results file path %q\n", resultsFilePath)

	resultFile, err := os.OpenFile(resultsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error al abrir/crear el archivo de resultados: %v", err)
	}
	defer resultFile.Close()

	// Función para procesar cada archivo
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Solo procesar archivos (no directorios)
		if !info.IsDir() {
			if !strings.Contains(path, resultsFileName) {
				err := searchInFile(path, searchTerm, resultFile)
				if err != nil {
					fmt.Printf("Error procesando el archivo %s: %v\n", path, err)
				}
			}

		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error al recorrer el directorio: %v", err)
	}

	return nil
}

// searchInFile busca el término en un archivo específico y guarda las líneas coincidentes en el archivo de resultados
func searchInFile(filePath, searchTerm string, resultFile *os.File) error {
	// Abrir el archivo que se va a leer
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Leer el archivo línea por línea
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Verificar si la línea contiene el término de búsqueda
		if strings.Contains(line, searchTerm) {

			currentTime := time.Now()
			now := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())
			// newLine:=fmt.Sprintf("%s; %q\n", filePath, line)
			newLine := fmt.Sprintf("%s; %s; %q\n", now, filePath, line)

			// Escribir la línea en el archivo de resultados
			if _, err := resultFile.WriteString(newLine); err != nil {
				return fmt.Errorf("error al escribir en el archivo de resultados: %v", err)
			}
		}
	}

	// Verificar si hubo errores al escanear
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error al leer el archivo: %v", err)
	}

	return nil
}

func main() {
	// Verificar si se pasaron los argumentos correctos
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run programa.go <directorio> <término_de_búsqueda>")
		return
	}

	// Obtener los argumentos de la línea de comandos
	directory := os.Args[1]
	searchTerm := os.Args[2]

	// Llamar a la función que busca en los archivos
	err := searchInFiles(directory, searchTerm)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Búsqueda completada con éxito.")
	}
}
