package baremetal

import (
	"yunion.io/x/onecloud/pkg/scheduler/algorithm/predicates"
	"yunion.io/x/onecloud/pkg/scheduler/core"
)

type CPUPredicate struct {
	BasePredicate
}

func (p *CPUPredicate) Name() string {
	return "baremetal_cpu"
}

func (p *CPUPredicate) Clone() core.FitPredicate {
	return &CPUPredicate{}
}

func (p *CPUPredicate) Execute(u *core.Unit, c core.Candidater) (bool, []core.PredicateFailureReason, error) {
	h := predicates.NewPredicateHelper(p, u, c)
	d := u.SchedData()

	freeCPUCount := h.GetInt64("FreeCPUCount", 0)
	reqCPUCount := int64(d.Ncpu)
	if freeCPUCount < reqCPUCount {
		totalCPUCount := h.GetInt64("CPUCount", 0)
		h.AppendInsufficientResourceError(reqCPUCount, totalCPUCount, freeCPUCount)
		h.SetCapacity(0)
	} else {
		if reqCPUCount/freeCPUCount != 1 {
			h.Exclude2("cpu", freeCPUCount, reqCPUCount)
		} else {
			h.SetCapacity(1)
		}
	}

	return h.GetResult()
}
