""" Este script eh um script chamador da biblioteca de NLP.
https://github.com/erickrf/nlpnet."""
import sys
import nlpnet


nlpnet.set_data_dir('/Users/emartins/nlpnet-data/')
TARGER = nlpnet.POSTagger()
print TARGER.tag(" ".join(sys.argv[1:]))

