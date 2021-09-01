package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/exceptions"
	"github.com/savannahghi/engagement/pkg/engagement/application/common/helpers"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	log "github.com/sirupsen/logrus"
)

func respondWithError(w http.ResponseWriter, code int, err error) {
	errMap := errorcodeutil.ErrorMap(err)
	errBytes, err := json.Marshal(errMap)
	if err != nil {
		errBytes = []byte(fmt.Sprintf("error: %s", err))
	}
	respondWithJSON(w, code, errBytes)
}

func respondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(payload)
	if err != nil {
		log.Printf(
			"unable to write payload `%s` to the http.ResponseWriter: %s",
			string(payload),
			err,
		)
	}
}

func getUIDFlavourAndIsAnonymous(r *http.Request) (*string, *feedlib.Flavour, *bool, error) {
	if r == nil {
		return nil, nil, nil, fmt.Errorf("nil request")
	}

	uid, err := getStringVar(r, "uid")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't get `uid` path var")
	}

	flavourStr, err := getStringVar(r, "flavour")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't get `flavour` path var: %w", err)
	}

	flavour := feedlib.Flavour(flavourStr)
	if !flavour.IsValid() {
		return nil, nil, nil, fmt.Errorf("`%s` is not a valid feed flavour", err)
	}

	isAnonymous, err := getStringVar(r, "isAnonymous")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't get `isAnonymous path var")
	}

	a, err := strconv.ParseBool(isAnonymous)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse to `isAnonymous` : %v ", err)
	}

	return &uid, &flavour, &a, nil

}

type patchItemFunc func(ctx context.Context, uid string, flavour feedlib.Flavour, itemID string) (*feedlib.Item, error)

func patchItem(
	ctx context.Context,
	patchFunc patchItemFunc,
	w http.ResponseWriter,
	r *http.Request,
) {
	itemID, err := getStringVar(r, "itemID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	item, err := patchFunc(addUIDToContext(ctx, *uid), *uid, *flavour, itemID)
	if err != nil {
		if errors.Is(err, exceptions.ErrNilFeedItem) {
			respondWithError(w, http.StatusNotFound, err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	marshalled, err := json.Marshal(item)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, marshalled)
}

type patchNudgeFunc func(ctx context.Context, uid string, flavour feedlib.Flavour, nudgeID string) (*feedlib.Nudge, error)

func patchNudge(
	ctx context.Context,
	patchFunc patchNudgeFunc,
	w http.ResponseWriter,
	r *http.Request,
) {
	nudgeID, err := getStringVar(r, "nudgeID")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	uid, flavour, _, err := getUIDFlavourAndIsAnonymous(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	element, err := patchFunc(addUIDToContext(ctx, *uid), *uid, *flavour, nudgeID)
	if err != nil {
		if errors.Is(err, exceptions.ErrNilNudge) {
			respondWithError(w, http.StatusNotFound, err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	marshalled, err := json.Marshal(element)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	respondWithJSON(w, http.StatusOK, marshalled)
}

func getOptionalBooleanFilterQueryParam(
	r *http.Request,
	paramName string,
) (*feedlib.BooleanFilter, error) {
	val := r.FormValue(paramName)
	if val == "" {
		return nil, nil // optional
	}

	boolFilter := feedlib.BooleanFilter(val)
	if !boolFilter.IsValid() {
		return nil, fmt.Errorf("optional bool: `%s` is not a valid boolean filter value", val)
	}

	return &boolFilter, nil
}

func getRequiredBooleanFilterQueryParam(
	r *http.Request,
	paramName string,
) (feedlib.BooleanFilter, error) {
	val := r.FormValue(paramName)
	if val == "" {
		return "", fmt.Errorf("required BooleanFilter `%s` not set", paramName)
	}

	boolFilter := feedlib.BooleanFilter(val)
	if !boolFilter.IsValid() {
		return "", fmt.Errorf("required bool: `%s` is not a valid boolean filter value", val)
	}

	return boolFilter, nil
}

func getOptionalStatusQueryParam(
	r *http.Request,
	paramName string,
) (*feedlib.Status, error) {
	val, err := getStringVar(r, paramName)
	if err != nil {
		return nil, nil // this is an optional param
	}

	status := feedlib.Status(val)
	if !status.IsValid() {
		return nil, fmt.Errorf("`%s` is not a valid status", val)
	}

	return &status, nil
}

func getOptionalVisibilityQueryParam(
	r *http.Request,
	paramName string,
) (*feedlib.Visibility, error) {
	val, err := getStringVar(r, paramName)
	if err != nil {
		return nil, nil // this is an optional param
	}

	visibility := feedlib.Visibility(val)
	if !visibility.IsValid() {
		return nil, fmt.Errorf("`%s` is not a valid visibility value", val)
	}

	return &visibility, nil
}

func getOptionalFilterParamsQueryParam(
	r *http.Request,
	paramName string,
) (*helpers.FilterParams, error) {
	// expect the filter params value to be JSON encoded
	val, err := getStringVar(r, paramName)
	if err != nil {
		return nil, nil // this is an optional param
	}

	filterParams := &helpers.FilterParams{}
	err = json.Unmarshal([]byte(val), filterParams)
	if err != nil {
		return nil, fmt.Errorf(
			"filter params should be a valid JSON representation of `helpers.FilterParams`. `%s` is not",
			val,
		)
	}

	return filterParams, nil
}

func getStringVar(r *http.Request, varName string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("can't get string var from a nil request")
	}
	pathVars := mux.Vars(r)
	pathVar, found := pathVars[varName]
	if !found {
		return "", fmt.Errorf("the request does not have a path var named `%s`", varName)
	}
	return pathVar, nil
}

// SchemaHandler ...
func SchemaHandler() (http.Handler, error) {
	f, err := pkger.Open("github.com/savannahghi/engagement:/static/")
	if err != nil {
		return nil, fmt.Errorf("can't open pkger schema dir: %w", err)
	}
	defer f.Close()

	return http.StripPrefix("/schema", http.FileServer(f)), nil
}

func addUIDToContext(ctx context.Context, uid string) context.Context {
	return context.WithValue(
		context.Background(),
		firebasetools.AuthTokenContextKey,
		&auth.Token{UID: uid},
	)
}
