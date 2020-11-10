// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package generated

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/plugin/federation/fedruntime"
)

func (ec *executionContext) __resolve__service(ctx context.Context) (fedruntime.Service, error) {
	if ec.DisableIntrospection {
		return fedruntime.Service{}, errors.New("federated introspection disabled")
	}

	var sdl []string

	for _, src := range sources {
		if src.BuiltIn {
			continue
		}
		sdl = append(sdl, src.Input)
	}

	return fedruntime.Service{
		SDL: strings.Join(sdl, "\n"),
	}, nil
}

func (ec *executionContext) __resolve_entities(ctx context.Context, representations []map[string]interface{}) ([]fedruntime.Entity, error) {
	list := []fedruntime.Entity{}
	for _, rep := range representations {
		typeName, ok := rep["__typename"].(string)
		if !ok {
			return nil, errors.New("__typename must be an existing string")
		}
		switch typeName {

		case "Action":
			id0, err := ec.unmarshalNString2string(ctx, rep["id"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "id"))
			}
			id1, err := ec.unmarshalNInt2int(ctx, rep["sequenceNumber"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "sequenceNumber"))
			}

			entity, err := ec.resolvers.Entity().FindActionByIDAndSequenceNumber(ctx,
				id0, id1)
			if err != nil {
				return nil, err
			}

			list = append(list, entity)

		case "Event":
			id0, err := ec.unmarshalNString2string(ctx, rep["id"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "id"))
			}

			entity, err := ec.resolvers.Entity().FindEventByID(ctx,
				id0)
			if err != nil {
				return nil, err
			}

			list = append(list, entity)

		case "Feed":
			id0, err := ec.unmarshalNString2string(ctx, rep["uid"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "uid"))
			}
			id1, err := ec.unmarshalNFlavour2gitlabᚗslade360emrᚗcomᚋgoᚋfeedᚋgraphᚋfeedᚐFlavour(ctx, rep["flavour"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "flavour"))
			}

			entity, err := ec.resolvers.Entity().FindFeedByUIDAndFlavour(ctx,
				id0, id1)
			if err != nil {
				return nil, err
			}

			list = append(list, entity)

		case "Item":
			id0, err := ec.unmarshalNString2string(ctx, rep["id"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "id"))
			}
			id1, err := ec.unmarshalNInt2int(ctx, rep["sequenceNumber"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "sequenceNumber"))
			}

			entity, err := ec.resolvers.Entity().FindItemByIDAndSequenceNumber(ctx,
				id0, id1)
			if err != nil {
				return nil, err
			}

			list = append(list, entity)

		case "Nudge":
			id0, err := ec.unmarshalNString2string(ctx, rep["id"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "id"))
			}
			id1, err := ec.unmarshalNInt2int(ctx, rep["sequenceNumber"])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Field %s undefined in schema.", "sequenceNumber"))
			}

			entity, err := ec.resolvers.Entity().FindNudgeByIDAndSequenceNumber(ctx,
				id0, id1)
			if err != nil {
				return nil, err
			}

			list = append(list, entity)

		default:
			return nil, errors.New("unknown type: " + typeName)
		}
	}
	return list, nil
}
