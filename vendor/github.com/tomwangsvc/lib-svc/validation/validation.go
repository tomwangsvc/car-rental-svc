package validation

import (
	"fmt"
	"strings"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
)

const (
	MetadataKeyIndex = lib_errors.MetadataKeyIndex
)

type Metadata map[string]interface{}

func ReduceMetadata(metadatas ...Metadata) Metadata {
	reduced := make(Metadata)
	for _, newMetadata := range metadatas {
		for newK, newV := range newMetadata {
			if newMDs, newOK := newV.([]Metadata); newOK {
				newV = reduceMetadataLists(nil, newMDs)
			}
			if reducedV, reducedOK := reduced[newK]; reducedOK {
				reducedMD, reducedOK := reducedV.(Metadata)
				newMD, newOK := newV.(Metadata)
				if reducedOK && newOK {
					reduced[newK] = ReduceMetadata(reducedMD, newMD)
					continue
				}
				reducedMDs, reducedOK := reducedV.([]Metadata)
				newMDs, newOK := newV.([]Metadata)
				if reducedOK && newOK {
					reduced[newK] = reduceMetadataLists(reducedMDs, newMDs)
					continue
				}
				reducedMessage, reducedOK := reducedV.(string) //TODO: consider erroring if not Metadata and not []Metadata and not string
				newMessage, newOK := newV.(string)             //TODO: consider erroring if not Metadata and not []Metadata and not string
				if reducedOK && newOK {
					reduced[newK] = ReduceMessages(Message(reducedMessage), Message(newMessage))
					continue
				}
				reduced[newK] = newV
				continue
			}
			reduced[newK] = newV
		}
	}
	return reduced
}

func reduceMetadataLists(reducedMDs, newMDs []Metadata) []Metadata {
	for _, newMD := range newMDs {
		if newIVal, newOK := newMD[MetadataKeyIndex]; newOK {
			var reduced bool
			for reducedI, reducedMD := range reducedMDs {
				if reducedIVal, reducedOK := reducedMD[MetadataKeyIndex]; reducedOK && newIVal == reducedIVal {
					reducedMDs[reducedI] = ReduceMetadata(reducedMD, newMD)
					reduced = true
					break
				}
			}
			if reduced {
				continue
			}
		}
		reducedMDs = append(reducedMDs, newMD)
	}
	return reducedMDs
}

type Message string

func ReduceMessages(messages ...Message) string {
	var ms string
	for _, m := range messages {
		if m := string(m); m != "" && !strings.Contains(ms, m) {
			if len(ms) != 0 {
				ms = fmt.Sprintf("%s,%s", ms, m)
			} else {
				ms = m
			}
		}
	}
	return ms
}
