// Model and express methods for a SCIM resource
// schemas and common attributes are represented explicitly
// core attributes and extension attributes are modelled as map
package resource

type Resource struct {
	Schemas    []string `json:"schemas"`
	Id         string   `json:"id"`
	ExternalId string   `json:"externalId"`
	Meta       *Meta    `json:"meta"`
	Attributes map[string]interface{}
}

// Create a new resource
func NewResource(resourceBaseUrn ...string) *Resource {
	return &Resource{
		Schemas:    resourceBaseUrn,
		Meta:       &Meta{},
		Attributes: make(map[string]interface{}, 0),
	}
}

// Create a new resource from map data
func NewResourceFromMap(data map[string]interface{}) *Resource {
	resource := NewResource()

	if schemas, ok := data["schemas"].([]string); ok {
		resource.Schemas = schemas
		delete(data, "schemas")
	}

	if id, ok := data["id"].(string); ok {
		resource.Id = id
		delete(data, "id")
	}

	if externalId, ok := data["externalId"].(string); ok {
		resource.ExternalId = externalId
		delete(data, "externalId")
	}

	if meta, ok := data["meta"].(map[string]interface{}); ok {
		if metaResourceType, ok := meta["resourceType"].(string); ok {
			resource.Meta.ResourceType = metaResourceType
		}
		if metaCreated, ok := meta["created"].(string); ok {
			resource.Meta.Created = metaCreated
		}
		if metaLastModified, ok := meta["lastModified"].(string); ok {
			resource.Meta.LastModified = metaLastModified
		}
		if metaLocation, ok := meta["location"].(string); ok {
			resource.Meta.Location = metaLocation
		}
		if metaVersion, ok := meta["version"].(string); ok {
			resource.Meta.Version = metaVersion
		}
		delete(data, "meta")
	}

	resource.Attributes = data
	return resource
}
