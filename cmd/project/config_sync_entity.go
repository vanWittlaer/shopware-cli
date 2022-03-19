package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"shopware-cli/shop"
)

type EntitySync struct{}

func (s EntitySync) Push(ctx adminSdk.ApiContext, client *adminSdk.Client, config *shop.Config, operation *ConfigSyncOperation) error {
	for _, entity := range config.Sync.Entity {
		if entity.Exists != nil && len(*entity.Exists) > 0 {
			criteria := make(map[string]interface{})
			criteria["filter"] = entity.Exists

			searchPayload, err := json.Marshal(criteria)

			if err != nil {
				return err
			}

			r, err := client.NewRequest(ctx, "POST", fmt.Sprintf("/api/search-ids/%s", entity.Entity), bytes.NewReader(searchPayload))

			if err != nil {
				return err
			}

			r.Header.Set("Accept", "application/json")
			r.Header.Set("Content-Type", "application/json")

			var res criteriaApiResponse
			_, err = client.Do(ctx.Context, r, &res)

			if err != nil {
				return err
			}

			if res.Total > 0 {
				continue
			}
		}

		operation.Operations[shop.NewUuid()] = adminSdk.SyncOperation{
			Action:  "upsert",
			Entity:  entity.Entity,
			Payload: []map[string]interface{}{entity.Payload},
		}
	}

	return nil
}

func (s EntitySync) Pull(ctx adminSdk.ApiContext, client *adminSdk.Client, config *shop.Config) error {
	return nil
}

type criteriaApiResponse struct {
	Total int      `json:"total"`
	Data  []string `json:"data"`
}
