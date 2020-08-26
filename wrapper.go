package bviper

import (
	"time"

	"github.com/go-masonry/mortar/interfaces/cfg"
	"github.com/spf13/viper"
)

type viperWrapper struct {
	instance *viper.Viper
}

func (v *viperWrapper) Get(key string) cfg.Value {
	return &valueWrapper{
		key:          key,
		viperWrapper: v,
	}
}

func (v *viperWrapper) Set(key string, value interface{}) {
	v.instance.Set(key, value)
}

func (v *viperWrapper) Map() map[string]interface{} {
	return v.instance.AllSettings()
}

func (v *viperWrapper) Implementation() interface{} {
	return v.instance
}

type valueWrapper struct {
	key string
	*viperWrapper
}

func (v *valueWrapper) IsSet() bool {
	return v.instance.IsSet(v.key)
}

func (v *valueWrapper) Raw() interface{} {
	return v.instance.Get(v.key)
}

func (v *valueWrapper) Bool() bool {
	return v.instance.GetBool(v.key)
}

func (v *valueWrapper) Int() int {
	return v.instance.GetInt(v.key)
}

func (v *valueWrapper) Int32() int32 {
	return v.instance.GetInt32(v.key)
}

func (v *valueWrapper) Int64() int64 {
	return v.instance.GetInt64(v.key)
}

func (v *valueWrapper) Uint() uint {
	return v.instance.GetUint(v.key)
}

func (v *valueWrapper) Uint32() uint32 {
	return v.instance.GetUint32(v.key)
}

func (v *valueWrapper) Uint64() uint64 {
	return v.instance.GetUint64(v.key)
}

func (v *valueWrapper) Float64() float64 {
	return v.instance.GetFloat64(v.key)
}

func (v *valueWrapper) Time() time.Time {
	return v.instance.GetTime(v.key)
}

func (v *valueWrapper) Duration() time.Duration {
	return v.instance.GetDuration(v.key)
}

func (v *valueWrapper) String() string {
	return v.instance.GetString(v.key)
}

func (v *valueWrapper) IntSlice() []int {
	return v.instance.GetIntSlice(v.key)
}

func (v *valueWrapper) StringSlice() []string {
	return v.instance.GetStringSlice(v.key)
}

func (v *valueWrapper) StringMap() map[string]interface{} {
	return v.instance.GetStringMap(v.key)
}

func (v *valueWrapper) StringMapString() map[string]string {
	return v.instance.GetStringMapString(v.key)
}

func (v *valueWrapper) StringMapStringSlice() map[string][]string {
	return v.instance.GetStringMapStringSlice(v.key)
}

// Unmarshal uses default decoder options, if you need some special behavior than it's best to get cfg.Implementation() and use it from there
func (v *valueWrapper) Unmarshal(result interface{}) error {
	return v.instance.UnmarshalKey(v.key, result)
}

var _ cfg.Config = (*viperWrapper)(nil)
var _ cfg.Value = (*valueWrapper)(nil)
