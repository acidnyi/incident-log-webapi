package incident_log

import (
	"slices"
)

func (l *IncidentLogDefinition) reconcileIncidents() {
	slices.SortFunc(l.Incidents, func(left, right Incident) int {
		if left.OccurredAt.Before(right.OccurredAt) {
			return -1
		} else if left.OccurredAt.After(right.OccurredAt) {
			return 1
		} else {
			return 0
		}
	})
}
