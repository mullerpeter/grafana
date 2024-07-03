package timeinterval

import (
	"encoding/json"
	"fmt"
	"hash/fnv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	model "github.com/grafana/grafana/pkg/apis/alerting_notifications/v0alpha1"
	"github.com/grafana/grafana/pkg/services/apiserver/endpoints/request"
	"github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
)

func getIntervalUID(t definitions.MuteTimeInterval) string {
	sum := fnv.New64()
	_, _ = sum.Write([]byte(t.Name))
	return fmt.Sprintf("%016x", sum.Sum64())
}

func convertToK8sResources(orgID int64, intervals []definitions.MuteTimeInterval, namespacer request.NamespaceMapper) (*model.TimeIntervalList, error) {
	data, err := json.Marshal(intervals)
	if err != nil {
		return nil, err
	}
	var specs []model.TimeIntervalSpec
	err = json.Unmarshal(data, &specs)
	if err != nil {
		return nil, err
	}
	result := &model.TimeIntervalList{}

	for idx := range specs {
		interval := intervals[idx]
		spec := specs[idx]
		item := buildTimeInterval(orgID, interval, spec, namespacer)
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func convertToK8sResource(orgID int64, interval definitions.MuteTimeInterval, namespacer request.NamespaceMapper) (*model.TimeInterval, error) {
	data, err := json.Marshal(interval)
	if err != nil {
		return nil, err
	}
	spec := model.TimeIntervalSpec{}
	err = json.Unmarshal(data, &spec)
	if err != nil {
		return nil, err
	}
	result := buildTimeInterval(orgID, interval, spec, namespacer)
	return &result, nil
}

func buildTimeInterval(orgID int64, interval definitions.MuteTimeInterval, spec model.TimeIntervalSpec, namespacer request.NamespaceMapper) model.TimeInterval {
	uid := getIntervalUID(interval) // TODO replace to stable UID when we switch to normal storage
	i := model.TimeInterval{
		TypeMeta: resourceInfo.TypeMeta(),
		ObjectMeta: metav1.ObjectMeta{
			UID:             types.UID(uid), // TODO This is needed to make PATCH work
			Name:            uid,            // TODO replace to stable UID when we switch to normal storage
			Namespace:       namespacer(orgID),
			ResourceVersion: interval.Version,
		},
		Spec: spec,
	}
	i.SetProvenanceStatus(string(interval.Provenance))
	return i
}

func convertToDomainModel(interval *model.TimeInterval) (definitions.MuteTimeInterval, error) {
	b, err := json.Marshal(interval.Spec)
	if err != nil {
		return definitions.MuteTimeInterval{}, err
	}
	result := definitions.MuteTimeInterval{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return definitions.MuteTimeInterval{}, err
	}
	result.Version = interval.ResourceVersion
	err = result.Validate()
	if err != nil {
		return definitions.MuteTimeInterval{}, err
	}
	return result, nil
}
