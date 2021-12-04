// Defines TensorFlow NLP modeler

package nlp

import (
	"fmt"
	"github.com/galeone/tfgo"
	tf "github.com/galeone/tensorflow/tensorflow/go"
	"github.com/icwells/dbIO"
)

type Processor struct {
	db		*dbIO.DBIO
	Model	*tfgo.Model
	root	*tf.Scope

}

func newProcessor(db *dbIO.DBIO) *Processor {
	p := new(Processor)
	p.db = db
	p.root := tfgo.NewRoot()
	return p
}

func NewNameClassifier(db *dbIO.DBIO) *Processor {
	// Returns new common/scientific name classifier
	p := newProcessor(db)
	
	return p
}

//func NewFuzzyMatcher(db *dbIO.DBIO) {}
