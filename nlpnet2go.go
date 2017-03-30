package main

import (
	"encoding/json"
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
		//fmt.Printf("%s\n", string(cmdReader))
		/*
			Os replaces abaixo são para tratar o retorno do nlpnet.
			Exemplo: [[(u'o', u'ART'), (u'rato', u'N'), (u'roeu', u'V'), (u'a', u'ART'), (u'roupa', u'N'), (u'do', u'PREP+ART'), (u'rei', u'N'), (u'de', u'PREP'), (u'roma', u'N'), (u'com', u'PREP'), (u'queijo', u'N')]]
		*/
		replacer := strings.NewReplacer("[", "", "]", "", "(u'", "", " u'", "", "'", "", "),", ");")
		replacerInsider := strings.NewReplacer("(", "", ")", "")
		result = replacer.Replace(string(cmdReader))
		retornos := strings.Split(result, ";")
		result = ""
		ret := Retorno{}
		arrRetornoItems := []RetornoItem{}
		for _, item := range retornos {
			item = replacerInsider.Replace(item)
			chaveValores := strings.Split(item, ",")
			tipo := strings.TrimSpace(chaveValores[1])
			valor := strings.TrimSpace(chaveValores[0])
			if tipo == "N" || tipo == "NPROP" {
				if valor == "queijo" {
					tipo = "ALIMENTO"
				}
			}
			retItem := RetornoItem{Chave: tipo, Valor: valor}
			arrRetornoItems = append(arrRetornoItems, retItem)
			fmt.Printf("Tipo: [%s] Valor: [%s]\n", tipo, valor)
		}
		ret.RetornoItems = arrRetornoItems
		jsonret, err := json.Marshal(ret)
		if err != nil {
			return "Erro na geracao do JSON. Erro: " + err.Error()
		}
		return string(jsonret)
	} else {
		return string(cmdReader)
	}

}

type attributes struct {
	Method      string `ini:"method"`
	Port        string `ini:"port"`
	CmdName     string `ini:"cmdName"`
	Nlpnet2gopy string `ini:"nlpnet2gopy"`
	Setdatadir  string `ini:"setdatadir"`
}

type Retorno struct {
	RetornoItems []RetornoItem `json:"items"`
}
type RetornoItem struct {
	Chave string `json:"chave"`
	Valor string `json:"valor"`
}
