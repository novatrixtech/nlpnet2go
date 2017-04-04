# nlpnet2go
Go Language wrapper and a web service for NLP Net project in Python - https://github.com/erickrf/nlpnet . NLP is a neural network architecture for NLP tasks, inspired in the SENNA system. It uses as input vectors built through Vector Space Models and avoid external NLP tools. Currently, it can perform POS tagging, SRL and dependency parsing. It has already processed models in Portuguese available at: http://nilc.icmc.usp.br/nlpnet/models.html 

If you need a POS (part-of-speech) analysis in Portuguese, this is your tool. 

Em bom português: **se tu precisas de uma ferramenta de processamento de linguagem natural para fazer a analise de parte de texto EM PORTUGUÊS, esta é a tua ferramenta.** As outras são muito complicadas. Acredite.

NLP net was created by Erick Rocha Fonseca - https://github.com/erickrf

To use it you need to setup NLP net first than just point out to nlpnet2go system where to call NLP net. 
In order to install NLP Net please follow the steps below:
```
$ python -m pip install --upgrade pip
$ pip install --user numpy scipy matplotlib ipython jupyter pandas sympy nose nltk nlpnet 
$ python
>>> import nltk
>>> nltk.download()
```
Go to the Models tab in the NLTK screen and select the Punkt tokenizer. It is used in order to split the text into sentences.
