package hera

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

func mergeStringMap(dest, src map[string]interface{}) {
	for sk, sv := range src {
		tv, ok := dest[sk]
		if !ok {
			dest[sk] = sv
			continue
		}

		svType := reflect.TypeOf(sv)
		tvType := reflect.TypeOf(tv)
		if svType != tvType {
			fmt.Println("continue, type is different")
			continue
		}

		switch ttv := tv.(type) {
		case map[interface{}]interface{}:
			tsv := sv.(map[interface{}]interface{})
			ssv := toMapStringInterface(tsv)
			stv := toMapStringInterface(ttv)
			mergeStringMap(ssv, stv)
			dest[sk] = ssv
		case map[string]interface{}:
			mergeStringMap(sv.(map[string]interface{}), ttv)
			dest[sk] = sv
		default:
			dest[sk] = sv
		}
	}
}

func toMapStringInterface(src map[interface{}]interface{}) map[string]interface{} {
	tgt := map[string]interface{}{}
	for k, v := range src {
		tgt[fmt.Sprintf("%v", k)] = v
	}
	return tgt
}

func insensitiviseMap(m map[string]interface{}) {
	for key, val := range m {
		switch val.(type) {
		case map[interface{}]interface{}:
			insensitiviseMap(cast.ToStringMap(val))
		case map[string]interface{}:
			insensitiviseMap(val.(map[string]interface{}))
		}

		lower := strings.ToLower(key)
		if key != lower {
			delete(m, key)
		}
		m[lower] = val
	}
}

func deepSearchInMap(m map[string]interface{}, paths ...string) map[string]interface{} {
	mtmp := make(map[string]interface{})
	for k, v := range m {
		mtmp[k] = v
	}
	for _, k := range paths {
		m2, ok := mtmp[k]
		if !ok {
			m3 := make(map[string]interface{})
			mtmp[k] = m3
			mtmp = m3
			continue
		}

		m3, err := cast.ToStringMapE(m2)
		if err != nil {
			m3 = make(map[string]interface{})
			mtmp[k] = m3
		}
		// continue search
		mtmp = m3
	}
	return mtmp
}
