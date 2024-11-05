// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v0alpha1

// RouteApplyConfiguration represents a declarative configuration of the Route type for use
// with apply.
type RouteApplyConfiguration struct {
	Continue          *bool                       `json:"continue,omitempty"`
	GroupBy           []string                    `json:"group_by,omitempty"`
	GroupInterval     *string                     `json:"group_interval,omitempty"`
	GroupWait         *string                     `json:"group_wait,omitempty"`
	Matchers          []MatcherApplyConfiguration `json:"matchers,omitempty"`
	MuteTimeIntervals []string                    `json:"mute_time_intervals,omitempty"`
	Receiver          *string                     `json:"receiver,omitempty"`
	RepeatInterval    *string                     `json:"repeat_interval,omitempty"`
	Routes            []RouteApplyConfiguration   `json:"routes,omitempty"`
}

// RouteApplyConfiguration constructs a declarative configuration of the Route type for use with
// apply.
func Route() *RouteApplyConfiguration {
	return &RouteApplyConfiguration{}
}

// WithContinue sets the Continue field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Continue field is set to the value of the last call.
func (b *RouteApplyConfiguration) WithContinue(value bool) *RouteApplyConfiguration {
	b.Continue = &value
	return b
}

// WithGroupBy adds the given value to the GroupBy field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the GroupBy field.
func (b *RouteApplyConfiguration) WithGroupBy(values ...string) *RouteApplyConfiguration {
	for i := range values {
		b.GroupBy = append(b.GroupBy, values[i])
	}
	return b
}

// WithGroupInterval sets the GroupInterval field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GroupInterval field is set to the value of the last call.
func (b *RouteApplyConfiguration) WithGroupInterval(value string) *RouteApplyConfiguration {
	b.GroupInterval = &value
	return b
}

// WithGroupWait sets the GroupWait field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GroupWait field is set to the value of the last call.
func (b *RouteApplyConfiguration) WithGroupWait(value string) *RouteApplyConfiguration {
	b.GroupWait = &value
	return b
}

// WithMatchers adds the given value to the Matchers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Matchers field.
func (b *RouteApplyConfiguration) WithMatchers(values ...*MatcherApplyConfiguration) *RouteApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithMatchers")
		}
		b.Matchers = append(b.Matchers, *values[i])
	}
	return b
}

// WithMuteTimeIntervals adds the given value to the MuteTimeIntervals field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the MuteTimeIntervals field.
func (b *RouteApplyConfiguration) WithMuteTimeIntervals(values ...string) *RouteApplyConfiguration {
	for i := range values {
		b.MuteTimeIntervals = append(b.MuteTimeIntervals, values[i])
	}
	return b
}

// WithReceiver sets the Receiver field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Receiver field is set to the value of the last call.
func (b *RouteApplyConfiguration) WithReceiver(value string) *RouteApplyConfiguration {
	b.Receiver = &value
	return b
}

// WithRepeatInterval sets the RepeatInterval field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RepeatInterval field is set to the value of the last call.
func (b *RouteApplyConfiguration) WithRepeatInterval(value string) *RouteApplyConfiguration {
	b.RepeatInterval = &value
	return b
}

// WithRoutes adds the given value to the Routes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Routes field.
func (b *RouteApplyConfiguration) WithRoutes(values ...*RouteApplyConfiguration) *RouteApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithRoutes")
		}
		b.Routes = append(b.Routes, *values[i])
	}
	return b
}