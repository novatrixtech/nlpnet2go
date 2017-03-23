""" Este script eh um script chamador da biblioteca de NLP.
https://github.com/erickrf/nlpnet."""
import sys
try:
    from configparser import ConfigParser
except ImportError:
    from ConfigParser import ConfigParser  # ver. < 3.0
import nlpnet


CONFIG = ConfigParser()
CONFIG.read('setup.ini')
nlpnet.set_data_dir(CONFIG.get('attributes', 'setdatadir'))
METHOD = CONFIG.get('attributes', 'method')
if METHOD == "pos":
    TAGGER = nlpnet.POSTagger()
    print TAGGER.tag(" ".join(sys.argv[1:]))
elif METHOD == "srl":
    TAGGER = nlpnet.SRLTagger()
    SENT = TAGGER.tag(" ".join(sys.argv[1:]))[0]
    print SENT.arg_structures
else:
    print TAGGER.tag("Invalid Tagger method operator. Only 'pos' OR 'srl' allowed.")


