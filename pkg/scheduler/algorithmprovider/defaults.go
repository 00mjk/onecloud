package algorithmprovider

import (
	"yunion.io/x/pkg/util/sets"

	"yunion.io/x/onecloud/pkg/scheduler/algorithm/predicates"
	predicateguest "yunion.io/x/onecloud/pkg/scheduler/algorithm/predicates/guest"
	priorityguest "yunion.io/x/onecloud/pkg/scheduler/algorithm/priorities/guest"
	"yunion.io/x/onecloud/pkg/scheduler/factory"
)

func init() {
	factory.RegisterAlgorithmProvider(factory.DefaultProvider, defaultPredicates(), defaultPriorities())
}

func defaultPredicates() sets.String {
	return sets.NewString(
		factory.RegisterFitPredicate("a-GuestHostStatusFilter", &predicateguest.StatusPredicate{}),
		factory.RegisterFitPredicate("b-GuestHypervisorFilter", &predicateguest.HypervisorPredicate{}),
		factory.RegisterFitPredicate("c-GuestAggregateFilter", &predicates.AggregatePredicate{}),
		factory.RegisterFitPredicate("d-GuestMigrateFilter", &predicateguest.MigratePredicate{}),
		factory.RegisterFitPredicate("e-GuestNestFilter", &predicateguest.NestPredicate{}),
		factory.RegisterFitPredicate("f-GuestGroupFilter", &predicateguest.GroupPredicate{}),
		factory.RegisterFitPredicate("g-GuestCPUFilter", &predicateguest.CPUPredicate{}),
		factory.RegisterFitPredicate("h-GuestMemoryFilter", &predicateguest.MemoryPredicate{}),
		factory.RegisterFitPredicate("i-GuestStorageFilter", &predicateguest.StoragePredicate{}),
		factory.RegisterFitPredicate("j-GuestNetworkFilter", &predicateguest.NetworkPredicate{}),
		factory.RegisterFitPredicate("k-GuestIsolatedDeviceFilter", &predicateguest.IsolatedDevicePredicate{}),
		factory.RegisterFitPredicate("l-GuestResourceTypeFilter", &predicates.ResourceTypePredicate{}),
	)
}

func defaultPriorities() sets.String {
	return sets.NewString(
		factory.RegisterPriority("guest-avoid-same-host", &priorityguest.AvoidSameHostPriority{}, 1),
		factory.RegisterPriority("guest-lowload", &priorityguest.LowLoadPriority{}, 1),
		factory.RegisterPriority("guest-creating", &priorityguest.CreatingPriority{}, 1),
		factory.RegisterPriority("guest-capacity", &priorityguest.CapacityPriority{}, 1),
	)
}
