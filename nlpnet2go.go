package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {

	args := os.Args[1:]

	//TODO: Carregar setup através do .INI
	attribts := attributes{}
	attribts.cmdName = "python"
	attribts.nlpnet2gopy = "nlpnet2go.py"
	attribts.method = "pos"
	attribts.port = ":5000"

	/*Caso existam argumentos de entrada,
	se comporta apenas como aplicação de console,
	caso contrário se comporta como serviço de api web*/
	if len(args) != 0 {
		for i := 0; i < len(args); i++ {
			attribts.text2analize += args[i] + " "
		}
		log.Println(execNlpnetCommand(attribts))

	} else {
		http.HandleFunc("/parser", func(w http.ResponseWriter, r *http.Request) {

			/* Valida as entradas e carrega a estrutura */
			if method := r.FormValue("method"); method != "" {
				attribts.method = method
			}
			if txt := r.FormValue("txt"); txt != "" {
				attribts.text2analize = txt
			}

			/* Retorna a mensagem analisada */
			io.WriteString(w, execNlpnetCommand(attribts))

		})

		http.ListenAndServe(attribts.port, nil)
	}

}

func execNlpnetCommand(attribts attributes) string {

	cmdArgs := []string{attribts.nlpnet2gopy, attribts.text2analize}
	cmd := exec.Command(attribts.cmdName, cmdArgs...)

	/* Opção A código minificado */
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// log.Println(cmd.Run())

	/* Opção B código tratado */
	cmdReader, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	//TODO: Converter JSON
	replacer := strings.NewReplacer("(u'", "{\"", "', u'", "\": \"", "')", "\"}")
	//, "',", "\":", "(", "{", ")", "}", "'", "\"", "\",", "\":")

	return replacer.Replace(string(cmdReader))
}

type attributes struct {
	method       string
	port         string
	text2analize string
	cmdName      string
	nlpnet2gopy  string
}
