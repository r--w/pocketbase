package pocketbase

import (
	"testing"
	"time"

	"github.com/r--w/pocketbase/migrations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollection_List(t *testing.T) {
	defaultClient := NewClient(defaultURL)

	tests := []struct {
		name       string
		client     *Client
		collection string
		params     ParamsList
		wantResult bool
		wantErr    bool
	}{
		{
			name:       "List with no params",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			wantErr:    false,
			wantResult: true,
		},
		{
			name:       "List no results - query",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			params: ParamsList{
				Filters: "field='some_random_value'",
			},
			wantErr:    false,
			wantResult: false,
		},
		{
			name:       "List no results - invalid query",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			params: ParamsList{
				Filters: "field~~~some_random_value'",
			},
			wantErr:    true,
			wantResult: false,
		},
		{
			name:       "List invalid collection",
			client:     defaultClient,
			collection: "invalid_collection",
			wantErr:    true,
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection := Collection[map[string]any]{tt.client, tt.collection}
			got, err := collection.List(tt.params)
			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.wantResult, got.TotalItems > 0)
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	client := NewClient(defaultURL)
	field := "value_" + time.Now().Format(time.StampMilli)
	collection := Collection[map[string]any]{client, migrations.PostsPublic}

	// delete non-existing item
	err := collection.Delete("non_existing_id")
	assert.Error(t, err)

	// create temporary item
	resultCreated, err := collection.Create(map[string]any{
		"field": field,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resultCreated.ID)

	// confirm item exists
	resultList, err := collection.List(ParamsList{Filters: "id='" + resultCreated.ID + "'"})
	assert.NoError(t, err)
	assert.Len(t, resultList.Items, 1)

	// delete temporary item
	err = collection.Delete(resultCreated.ID)
	assert.NoError(t, err)

	// confirm item does not exist
	resultList, err = collection.List(ParamsList{Filters: "id='" + resultCreated.ID + "'"})
	assert.NoError(t, err)
	assert.Len(t, resultList.Items, 0)
}

func TestCollection_Update(t *testing.T) {
	client := NewClient(defaultURL)
	field := "value_" + time.Now().Format(time.StampMilli)
	collection := Collection[map[string]any]{client, migrations.PostsPublic}

	// update non-existing item
	err := collection.Update("non_existing_id", map[string]any{
		"field": field,
	})
	assert.Error(t, err)

	// create temporary item
	resultCreated, err := collection.Create(map[string]any{
		"field": field,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resultCreated.ID)

	// confirm item exists
	resultList, err := collection.List(ParamsList{Filters: "id='" + resultCreated.ID + "'"})
	assert.NoError(t, err)
	require.Len(t, resultList.Items, 1)
	assert.Equal(t, field, resultList.Items[0]["field"])

	// update temporary item
	err = collection.Update(resultCreated.ID, map[string]any{
		"field": field + "_updated",
	})
	assert.NoError(t, err)

	// confirm changes
	resultList, err = collection.List(ParamsList{Filters: "id='" + resultCreated.ID + "'"})
	assert.NoError(t, err)
	require.Len(t, resultList.Items, 1)
	assert.Equal(t, field+"_updated", resultList.Items[0]["field"])
}

func TestCollection_Create(t *testing.T) {
	defaultClient := NewClient(defaultURL)
	defaultBody := map[string]interface{}{
		"field": "value_" + time.Now().Format(time.StampMilli),
	}

	tests := []struct {
		name       string
		client     *Client
		collection string
		body       any
		wantErr    bool
		wantID     bool
	}{
		{
			name:       "Create with no body",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			wantErr:    true,
		},
		{
			name:       "Create with body",
			client:     defaultClient,
			collection: migrations.PostsPublic,
			body:       defaultBody,
			wantErr:    false,
			wantID:     true,
		},
		{
			name:       "Create invalid collections",
			client:     defaultClient,
			collection: "invalid_collection",
			body:       defaultBody,
			wantErr:    true,
		},
		{
			name:       "Create no auth",
			client:     defaultClient,
			collection: migrations.PostsUser,
			body:       defaultBody,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection := Collection[any]{tt.client, tt.collection}
			r, err := collection.Create(tt.body)
			assert.Equal(t, tt.wantErr, err != nil, err)
			assert.Equal(t, tt.wantID, r.ID != "")
		})
	}
}

func TestCollection_One(t *testing.T) {
	client := NewClient(defaultURL)
	field := "value_" + time.Now().Format(time.StampMilli)
	collection := Collection[map[string]any]{client, migrations.PostsPublic}

	// update non-existing item
	_, err := collection.One("non_existing_id")
	assert.Error(t, err)

	// create temporary item
	resultCreated, err := collection.Create(map[string]any{
		"field": field,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resultCreated.ID)

	// confirm item exists
	item, err := collection.One(resultCreated.ID)
	assert.NoError(t, err)
	assert.Equal(t, field, item["field"])

	// update temporary item
	err = collection.Update(resultCreated.ID, map[string]any{
		"field": field + "_updated",
	})
	assert.NoError(t, err)

	// confirm changes
	item, err = collection.One(resultCreated.ID)
	assert.NoError(t, err)
	assert.Equal(t, field+"_updated", item["field"])
}
