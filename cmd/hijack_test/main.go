package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Webserver doesn't support hijacking.", http.StatusInternalServerError)
			return
		}

		conn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer conn.Close()
		//buf.WriteString("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nsomebody once told me")
		conn.Write([]byte("Raw TCP time. Hello, client!"))
		return

		/*buf.WriteString("Raw TCP time. Say hi: ")
		buf.Flush()

		s, err := buf.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading string: %v\n", err)
			return
		}

		fmt.Fprintf(buf, "You said: %q\nBye.\n", s)
		buf.Flush()*/
	})

	http.ListenAndServe(":8000", nil)
}
