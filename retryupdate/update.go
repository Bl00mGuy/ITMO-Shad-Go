//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(client kvapi.Client, key string, updateFunction func(currentValue *string) (newValue string, err error)) error {
	initialValue, initialVersion, err := fetchInitialValue(client, key)
	if err != nil {
		return err
	}
	return attemptUpdateWithRetry(client, key, updateFunction, initialValue, initialVersion)
}

func fetchInitialValue(client kvapi.Client, key string) (*string, uuid.UUID, error) {
	var authError *kvapi.AuthError
	for {
		response, err := client.Get(&kvapi.GetRequest{Key: key})
		switch {
		case errors.As(err, &authError):
			return nil, uuid.UUID{}, err
		case err == nil:
			return &response.Value, response.Version, nil
		case errors.Is(err, kvapi.ErrKeyNotFound):
			return nil, uuid.UUID{}, nil
		}
	}
}

func attemptUpdateWithRetry(client kvapi.Client, key string, updateFunction func(currentValue *string) (newValue string, err error), currentValue *string, currentVersion uuid.UUID) error {
	var authError *kvapi.AuthError
	var conflictError *kvapi.ConflictError
	newVersionID := uuid.Must(uuid.NewV4())
	for {
		updatedValue, err := updateFunction(currentValue)
		if err != nil {
			return err
		}
		err = trySetUpdatedValue(client, key, updatedValue, currentVersion, newVersionID)
		switch {
		case errors.As(err, &authError):
			return err
		case err == nil:
			return nil
		case errors.Is(err, kvapi.ErrKeyNotFound):
			currentValue = nil
			currentVersion = uuid.UUID{}
			continue
		case errors.As(err, &conflictError):
			if conflictError.ExpectedVersion == newVersionID {
				return nil
			}
			return UpdateValue(client, key, updateFunction)
		}
	}
}

func trySetUpdatedValue(client kvapi.Client, key, updatedValue string, oldVersion, newVersion uuid.UUID) error {
	_, err := client.Set(&kvapi.SetRequest{
		Key:        key,
		Value:      updatedValue,
		OldVersion: oldVersion,
		NewVersion: newVersion,
	})
	return err
}
