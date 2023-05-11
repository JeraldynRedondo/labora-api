/*
Sirve para manejar la configuración de la aplicación, como la conexión a la base de datos.
*/

package config

import (
	"fmt"
	"net/http"
	"time"
)

func StartServer(port string, router http.Handler) error {
	servidor := &http.Server{
		Handler:      router,
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Starting Server on port %s...\n", port)
	if err := servidor.ListenAndServe(); err != nil {
		return fmt.Errorf("Error while starting up Server: '%v'", err)
	}
	return nil
}
