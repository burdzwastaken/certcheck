package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/certifi/gocertifi"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	host := r.Form.Get("text") + ":443"

	certPool, err := gocertifi.CACerts()

	tlsconfig := &tls.Config{
		RootCAs: certPool,
	}

	conn, err := tls.Dial("tcp", host, tlsconfig)
	if err != nil {
		fmt.Println("err:", err)
	}

	defer conn.Close()

	for _, chain := range conn.ConnectionState().VerifiedChains {
		for i := len(chain) - 3; i >= 0; i-- {
			cert := chain[i]
			fmt.Fprintf(w, "certificate name: %s \ncertificate SANs: %s \n", cert.Subject.CommonName, cert.DNSNames)
		}
	}

}
