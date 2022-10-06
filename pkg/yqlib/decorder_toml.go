package yqlib

import (
	"io"

	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v3"
)

type tomlDecoder struct {
	reader       io.Reader
	readAnything bool
	finished     bool
}

func NewTomlDecoder() Decoder {
	return &tomlDecoder{}
}

func (dec *tomlDecoder) Init(reader io.Reader) {
	dec.reader = reader
	dec.readAnything = false
	dec.finished = false
}

func (dec *tomlDecoder) keyDiff(context toml.Key, key toml.Key) (toml.Key, bool) {
	if len(context) >= len(key) {
		log.Debug("keyDiff Finished: new context? %v vs %v", context, key)
		return nil, true
	}

	//check keys are a subset of context
	for i, contextValue := range context {
		if key[i] != contextValue {
			log.Debug("keyDiff Finished: new context? %v vs %v", context, key)
			return nil, true
		}
	}
	//key is a subset, return the extra bits
	log.Debug("keyDiff result: new context? %v", key[len(context):])
	return key[len(context):], false
}

func (dec *tomlDecoder) sortMap(context toml.Key, keys []toml.Key, i int, data *orderedMap) int {
	log.Debugf("sortMap: need to sort this map %v", data)
	j := i
	newKv := make([]orderedMapKV, len(data.kv))
	newKvIndex := 0
	for {
		if j == len(keys) {
			break
		}
		diff, finished := dec.keyDiff(context, keys[j])

		if finished || len(diff) != 1 {
			log.Debug("sortMap: finished, diff: %v, keys: %v, j: %v", diff, keys[j], j)
			break
		}
		log.Debug("sortMap: sorting %v into position %v", diff[0], newKvIndex)
		newKv[newKvIndex].K = diff[0]
		newKv[newKvIndex].V = data.GetKey(diff[0])

		newKvIndex = newKvIndex + 1
		j = j + 1
	}

	data.kv = newKv
	log.Debug("sortMap: result: %v", newKv)
	return j
}

func (dec *tomlDecoder) reorderMaps(currentContext toml.Key, keys []toml.Key, i int, data *orderedMap) int {
	// context: ["a", "b"], meaning I'm in "b" atm
	// keys: ["a", "b", "c"], meaning, I need to sort "c" of "b", the children of "b"

	if i == len(keys) {
		log.Debug("reorderMaps - done at %v", i)
		return i
	}
	log.Debug("reorderMaps, %v, %v - %v", currentContext, keys[i], i)

	diff, finished := dec.keyDiff(currentContext, keys[i])

	if finished {
		return i
	} else if len(diff) == 1 { // we need to sort this map
		j := dec.sortMap(currentContext, keys, i, data)
		log.Debug("sortMap results here:", data)
		return dec.reorderMaps(currentContext, keys, j, data)
	} // we need to sort a child map
	childKey := keys[i][0]
	newContext := append(currentContext, childKey)
	child := data.GetKey(childKey)
	j := dec.reorderMaps(newContext, keys, i, &child)
	log.Debug("child results here:", child)
	// process the next item
	return dec.reorderMaps(currentContext, keys, j, data)

}

func (dec *tomlDecoder) Decode(rootYamlNode *yaml.Node) error {
	if dec.finished {
		return io.EOF
	}
	var data orderedMap
	decoder := toml.NewDecoder(dec.reader)
	metadata, err := decoder.Decode(&data)

	if err != nil {
		return err
	}

	// decoder does not maintain key order :( have to iterate through the structure and order
	// all the map keys

	log.Debugf("metadata: %v", metadata)
	log.Debugf("metadata: %v", metadata.Keys())
	dec.reorderMaps(toml.Key{}, metadata.Keys(), 0, &data)
	log.Debug("reoreder results here:", data)

	dec.finished = true
	node, err := data.convertToYamlNode()
	if err != nil {
		return err
	}
	rootYamlNode.Kind = yaml.DocumentNode
	rootYamlNode.Content = []*yaml.Node{node}
	return nil
}
