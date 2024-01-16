package main

import (
	"io"
	"net/http"
	"os"
	"serversidemapper32/ssm32"
)

// main is the entry point of the program.
func main() {
	// Read the content of the "helloworld32.dll" file.
	dllContent, err := os.ReadFile("helloworld32.dll")
	if err != nil {
		panic(err)
	}

	// Get the manual map data and context for the DLL.
	dllMMap32Data, dllMMap32Context := ssm32.GetMMap32Data(dllContent)

	// Handle HTTP requests.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/get_mmap_data":
			// Serve the manual map data.
			w.Write(dllMMap32Data)
		case "/get_mapped_dll":
			// Process the request body.
			processedData, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Map the DLL using the processed data and manual map context.
			mappedDll := ssm32.MMap32(processedData, dllMMap32Context)
			if mappedDll == nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Write(mappedDll)
			}
		}
	})

	// Start the HTTP server on port 8000.
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
