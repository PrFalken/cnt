package builder

import (
	"github.com/blablacar/cnt/log"
)

func (p *Pod) Test() {
	log.Get().Info("Testing POD", p.manifest.Name)

	for _, e := range p.manifest.Pod.Apps {
		aci, err := NewAciWithManifest(p.path+"/"+e.Name, p.args, p.toAciManifest(e))
		if err != nil {
			log.Get().Panic(err)
		}
		aci.PodName = &p.manifest.Name
		aci.Test()
	}
}
