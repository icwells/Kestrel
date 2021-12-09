// Defines TensorFlow NLP modeler

package terms

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"github.com/galeone/tfgo"
	tf "github.com/galeone/tensorflow/tensorflow/go"
)

type Processor struct {
	infile	string
	model	*tfgo.Model
}

func newProcessor(db *dbIO.DBIO) *Processor {
	p := new(Processor)
	p.infile = path.Join(iotools.GetGOPATH(), "src/github.com/icwells/kestrel/nlp/nlpModel")
	p.Model = tfgo.LoadModel(p.infile, []string{"serve"}, nil)
	return p
}

func (p *Processor) Common(val string) bool {
	// Returns true if a name is common
	



	return ret
}
