package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/go-ini/ini"
)

func main() {

	args := os.Args[1:]
	text2analize := ""
	attribts := new(attributes)
	setupfile := "setup.ini"

	//Carregar setup através do .ini
	cfg, err := ini.Load(setupfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	err = cfg.Section("attributes").MapTo(attribts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	/*Caso existam argumentos de entrada,
	se comporta apenas como aplicação de console,
	caso contrário se comporta como serviço de api web*/
	if len(args) != 0 {
		for i := 0; i < len(args); i++ {
			if args[i] == "-m" {
				attribts.Method = args[i+1]
				i = i + 1
			} else {
				text2analize += args[i] + " "
			}

		}
		log.Println(execNlpnetCommand(attribts, text2analize))

	} else {
		http.HandleFunc("/parser", func(w http.ResponseWriter, r *http.Request) {

			/* Valida as entradas e carrega a estrutura */
			if method := r.FormValue("method"); method != "" {
				attribts.Method = method
			}
			if txt := r.FormValue("txt"); txt != "" {
				text2analize = txt
			}

			/* Retorna a mensagem analisada */
			io.WriteString(w, execNlpnetCommand(attribts, text2analize))

		})
		fmt.Printf("nlpnet2go running at: %s\n", attribts.Port)
		http.ListenAndServe(":"+attribts.Port, nil)

	}

}

func execNlpnetCommand(attribts *attributes, text2analize string) string {

	cmdArgs := []string{attribts.Nlpnet2gopy, "-t", text2analize, "-m", attribts.Method}
	cmd := exec.Command(attribts.CmdName, cmdArgs...)

	cmdReader, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error: [%s]", err.Error())
		os.Exit(1)
	}

	/* Conversão para JSON */
	var result string
	if attribts.Method == "pos" {
		replacer := strings.NewReplacer("(u'", "{\"", "', u'", "\": \"", "')", "\"}")
		result = replacer.Replace(string(cmdReader))
	} else {
		result = string(cmdReader)
	}
	return result

}

type attributes struct {
	Method      string `ini:"method"`
	Port        string `ini:"port"`
	CmdName     string `ini:"cmdName"`
	Nlpnet2gopy string `ini:"nlpnet2gopy"`
	Setdatadir  string `ini:"setdatadir"`
}
