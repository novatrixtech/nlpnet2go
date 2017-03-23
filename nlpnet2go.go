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
			text2analize += args[i] + " "
		}
		log.Println(execNlpnetCommand(attribts, text2analize))

	} else {
		http.HandleFunc("/parser", func(w http.ResponseWriter, r *http.Request) {

			/* Valida as entradas e carrega a estrutura */
			if method := r.FormValue("method"); method != "" {
				attribts.Method = method

				/* Atualiza o arquivo de configuração em relação ao metodo de tagger
				a fim de ser utilizado pelo script python */
				cfg.Section("attributes").Key("method").SetValue(attribts.Method)
				err = cfg.SaveTo(setupfile)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err)
					os.Exit(1)
				}
			}
			if txt := r.FormValue("txt"); txt != "" {
				text2analize = txt
			}

			/* Retorna a mensagem analisada */
			io.WriteString(w, execNlpnetCommand(attribts, text2analize))

		})

		http.ListenAndServe(":"+attribts.Port, nil)
	}

}

func execNlpnetCommand(attribts *attributes, text2analize string) string {

	cmdArgs := []string{attribts.Nlpnet2gopy, text2analize}
	cmd := exec.Command(attribts.CmdName, cmdArgs...)

	cmdReader, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
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
