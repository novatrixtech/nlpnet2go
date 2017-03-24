""" Este script eh um script chamador da biblioteca de NLP.
https://github.com/erickrf/nlpnet."""
import sys
import getopt
try:
    from configparser import ConfigParser
except ImportError:
    from ConfigParser import ConfigParser  # ver. < 3.0
import nlpnet


CONFIG = ConfigParser()
CONFIG.read('setup.ini')
nlpnet.set_data_dir(CONFIG.get('attributes', 'setdatadir'))

TEXT = ''
METHOD = ''

try:
    OPTS, ARGS = getopt.getopt(sys.argv[1:], "ht:m:", ["text=", "method="])
except getopt.GetoptError:
    sys.exit(1)
for opt, arg in OPTS:
    if opt == '-h':
        print 'nlpnet2go.py -t <"text to be analyzed"> -m <method [''pos''] OR [''srl'']>'
        print 'Eg.: python nlpnet2go.py -t "teste do edward" -m pos'
        sys.exit()
    elif opt in ("-t", "--text"):
        TEXT = arg
    elif opt in ("-m", "--method"):
        METHOD = arg


if METHOD == "pos":
    TAGGER = nlpnet.POSTagger()
    print TAGGER.tag(TEXT)
elif METHOD == "srl":
    TAGGER = nlpnet.SRLTagger()
    SENT = TAGGER.tag(TEXT)[0]
    print SENT.arg_structures
else:
    print sys.argv[1:], "Invalid Tagger method operator. Only 'pos' OR 'srl' allowed."
# print sys.argv[1:]
