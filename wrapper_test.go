package bviper

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var timeStamp, _ = time.Parse(time.RFC3339, "2020-07-15T12:07:18+00:00")

type testStruct struct {
	One string `mapstructure:"one"`
}

func prepareViperWithDefaults() *viper.Viper {
	v := viper.New()
	v.SetDefault("bool", true)
	v.SetDefault("int", 1)
	v.SetDefault("int32", int32(32))
	v.SetDefault("int64", int64(64))
	v.SetDefault("uint", uint(8))
	v.SetDefault("uint32", uint32(32))
	v.SetDefault("uint64", uint64(64))
	v.SetDefault("float64", float64(64.0))
	v.SetDefault("time", timeStamp)
	v.SetDefault("duration", time.Hour)
	v.SetDefault("string", "string value")
	v.SetDefault("intslice", []int{1, 2, 3})
	v.SetDefault("stringslice", []string{"one", "two"})
	v.SetDefault("stringmap", map[string]interface{}{
		"one": 1,
	})
	v.SetDefault("stringmapstring", map[string]string{
		"one": "two",
	})
	v.SetDefault("stringmapstringslice", map[string][]string{
		"one": []string{"two", "three"},
	})
	v.SetDefault("object", map[string]interface{}{
		"one": "two",
	})
	return v
}

func TestWrapper(t *testing.T) {
	builder := CustomBuilder(prepareViperWithDefaults())
	config, err := builder.Build()
	assert.NoError(t, err)
	value := config.Get("some.key")
	assert.NotNil(t, value)
	impl := config.Implementation()
	assert.IsType(t, &viper.Viper{}, impl)
	keyMap := config.Map()
	assert.Contains(t, keyMap, "int")
	config.Set("some.key", map[string]interface{}{
		"name":  "viper",
		"sound": "ssss",
	})
	assert.True(t, config.Get("some.key").IsSet())
}

func TestValueWrapper(t *testing.T) {
	builder := CustomBuilder(prepareViperWithDefaults())
	config, err := builder.Build()
	assert.NoError(t, err)
	assert.False(t, config.Get("fake.key").IsSet())
	assert.True(t, config.Get("bool").IsSet())
	assert.IsType(t, true, config.Get("bool").Raw())
	assert.Equal(t, true, config.Get("bool").Bool())
	assert.Equal(t, 1, config.Get("int").Int())
	assert.Equal(t, int32(32), config.Get("int32").Int32())
	assert.Equal(t, int64(64), config.Get("int64").Int64())
	assert.Equal(t, uint(8), config.Get("uint").Uint())
	assert.Equal(t, uint32(32), config.Get("uint32").Uint32())
	assert.Equal(t, uint64(64), config.Get("uint64").Uint64())
	assert.Equal(t, float64(64.0), config.Get("float64").Float64())
	assert.Equal(t, timeStamp, config.Get("time").Time())
	assert.Equal(t, time.Hour, config.Get("duration").Duration())
	assert.Equal(t, "string value", config.Get("string").String())
	assert.Equal(t, []int{1, 2, 3}, config.Get("intslice").IntSlice())
	assert.Equal(t, []string{"one", "two"}, config.Get("stringslice").StringSlice())
	assert.Equal(t, map[string]interface{}{"one": 1}, config.Get("stringmap").StringMap())
	assert.Equal(t, map[string]string{"one": "two"}, config.Get("stringmapstring").StringMapString())
	assert.Equal(t, map[string][]string{"one": []string{"two", "three"}}, config.Get("stringmapstringslice").StringMapStringSlice())
	var obj testStruct
	err = config.Get("object").Unmarshal(&obj)
	assert.NoError(t, err)
	assert.EqualValues(t, obj, testStruct{One: "two"})
}

func TestWithFile(t *testing.T) {
	build, err := Builder().SetConfigFile("./config.yml").AddExtraConfigFile("./config_test.yml").Build()
	assert.NoError(t, err)
	duration := build.Get("self.duration").Duration()
	assert.Equal(t, time.Hour, duration)
	num := build.Get("self.int").Int()
	assert.Equal(t, 2, num)
}

func TestEnvironmentValues(t *testing.T) {
	err := os.Setenv("SOME_STRING", "value")
	assert.NoError(t, err)
	conf, err := Builder().Build()
	assert.NoError(t, err)
	value := conf.Get("some.string").String()
	assert.Equal(t, "value", value)
}

func TestEnvironmentDelimiterReplacer(t *testing.T) {
	err := os.Setenv("SOME__STRING", "value")
	assert.NoError(t, err)
	conf, err := Builder().SetEnvDelimiterReplacer("__", ".").Build()
	assert.NoError(t, err)
	value := conf.Get("some.string").String()
	assert.Equal(t, "value", value)
}
