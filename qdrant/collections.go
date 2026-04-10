package qdrant

import (
	"context"
	"errors"
)

// Checks the existence of a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - collectionName: The name of the collection to check.
//
// Returns:
//   - bool: True if the collection exists, false otherwise.
//   - error: An error if the operation fails.
func (c *Client) CollectionExists(ctx context.Context, collectionName string) (bool, error) {
	resp := &CollectionExistsResponse{}
	err := c.get().call(ctx, opCollectionExists, &CollectionExistsRequest{
		CollectionName: collectionName,
	}, resp)
	if err != nil {
		return false, newQdrantErr(err, "CollectionExists", collectionName)
	}
	return resp.GetResult().GetExists(), nil
}

// Retrieves detailed information about a specified existing collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - collectionName: The name of the collection to retrieve information for.
//
// Returns:
//   - *CollectionInfo: Detailed information about the collection.
//   - error: An error if the operation fails.
func (c *Client) GetCollectionInfo(ctx context.Context, collectionName string) (*CollectionInfo, error) {
	resp := &GetCollectionInfoResponse{}
	err := c.get().call(ctx, opGetCollection, &GetCollectionInfoRequest{
		CollectionName: collectionName,
	}, resp)
	if err != nil {
		return nil, newQdrantErr(err, "GetCollection", collectionName)
	}
	return resp.GetResult(), nil
}

// Retrieves the names of all existing collections.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - []string: A slice of collection names.
//   - error: An error if the operation fails.
func (c *Client) ListCollections(ctx context.Context) ([]string, error) {
	resp := &ListCollectionsResponse{}
	err := c.get().call(ctx, opListCollections, &ListCollectionsRequest{}, resp)
	if err != nil {
		return nil, newQdrantErr(err, "ListCollections")
	}
	var collections []string
	for _, collection := range resp.GetCollections() {
		collections = append(collections, collection.GetName())
	}
	return collections, nil
}

// Creates a new collection with the given parameters.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The CreateCollection request containing collection parameters.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) CreateCollection(ctx context.Context, request *CreateCollection) error {
	resp := &CollectionOperationResponse{}
	err := c.get().call(ctx, opCreateCollection, request, resp)
	if err != nil {
		return newQdrantErr(err, "CreateCollection", request.GetCollectionName())
	}
	return nil
}

// Updates parameters of an existing collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The UpdateCollection request containing updated parameters.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) UpdateCollection(ctx context.Context, request *UpdateCollection) error {
	resp := &CollectionOperationResponse{}
	err := c.get().call(ctx, opUpdateCollection, request, resp)
	if err != nil {
		return newQdrantErr(err, "UpdateCollection", request.GetCollectionName())
	}
	return nil
}

// Drops a collection and all associated data.
//
// Parameters:
//   - ctx: The context for the request.
//   - collectionName: The name of the collection to delete.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) DeleteCollection(ctx context.Context, collectionName string) error {
	resp := &CollectionOperationResponse{}
	err := c.get().call(ctx, opDeleteCollection, &DeleteCollection{
		CollectionName: collectionName,
	}, resp)
	if err != nil {
		return newQdrantErr(err, "DeleteCollection", collectionName)
	}
	if !resp.GetResult() {
		return newQdrantErr(errors.New("failed to delete collection"), "DeleteCollection", collectionName)
	}
	return nil
}

// Creates an alias for a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - aliasName: The name of the alias to create.
//   - collectionName: The name of the collection to alias.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) CreateAlias(ctx context.Context, aliasName, collectionName string) error {
	resp := &CollectionOperationResponse{}
	err := c.get().call(ctx, opUpdateAliases, &ChangeAliases{
		Actions: []*AliasOperations{
			{
				Action: &AliasOperations_CreateAlias{
					CreateAlias: &CreateAlias{
						CollectionName: collectionName,
						AliasName:      aliasName,
					},
				},
			},
		},
	}, resp)
	if err != nil {
		return newQdrantErr(err, "CreateAlias", collectionName)
	}
	return nil
}

// Deletes an alias.
//
// Parameters:
//   - ctx: The context for the request.
//   - aliasName: The name of the alias to delete.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) DeleteAlias(ctx context.Context, aliasName string) error {
	resp := &CollectionOperationResponse{}
	err := c.get().call(ctx, opUpdateAliases, &ChangeAliases{
		Actions: []*AliasOperations{
			{
				Action: &AliasOperations_DeleteAlias{
					DeleteAlias: &DeleteAlias{
						AliasName: aliasName,
					},
				},
			},
		},
	}, resp)
	if err != nil {
		return newQdrantErr(err, "DeleteAlias", aliasName)
	}
	return nil
}

// Renames an alias.
//
// Parameters:
//   - ctx: The context for the request.
//   - oldAliasName: The current name of the alias.
//   - newAliasName: The new name for the alias.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) RenameAlias(ctx context.Context, oldAliasName, newAliasName string) error {
	resp := &CollectionOperationResponse{}
	err := c.get().call(ctx, opUpdateAliases, &ChangeAliases{
		Actions: []*AliasOperations{
			{
				Action: &AliasOperations_RenameAlias{
					RenameAlias: &RenameAlias{
						OldAliasName: oldAliasName,
						NewAliasName: newAliasName,
					},
				},
			},
		},
	}, resp)
	if err != nil {
		return newQdrantErr(err, "RenameAlias")
	}
	return nil
}

// Lists all aliases for a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - collectionName: The name of the collection to list aliases for.
//
// Returns:
//   - []string: A slice of alias names.
//   - error: An error if the operation fails.
func (c *Client) ListCollectionAliases(ctx context.Context, collectionName string) ([]string, error) {
	resp := &ListAliasesResponse{}
	err := c.get().call(ctx, opListCollectionAliases, &ListCollectionAliasesRequest{
		CollectionName: collectionName,
	}, resp)
	if err != nil {
		return nil, newQdrantErr(err, "ListCollectionAliases", collectionName)
	}
	var aliases []string
	for _, alias := range resp.GetAliases() {
		aliases = append(aliases, alias.GetAliasName())
	}
	return aliases, nil
}

// Lists all aliases.
//
// Parameters:
//   - ctx: The context for the request.
//
// Returns:
//   - []*AliasDescription: A slice of AliasDescription objects.
//   - error: An error if the operation fails.
func (c *Client) ListAliases(ctx context.Context) ([]*AliasDescription, error) {
	resp := &ListAliasesResponse{}
	err := c.get().call(ctx, opListAliases, &ListAliasesRequest{}, resp)
	if err != nil {
		return nil, newQdrantErr(err, "ListAliases")
	}
	return resp.GetAliases(), nil
}

// Updates aliases.
//
// Parameters:
//   - ctx: The context for the request.
//   - actions: A slice of AliasOperations to perform.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) UpdateAliases(ctx context.Context, actions []*AliasOperations) error {
	resp := &CollectionOperationResponse{}
	err := c.get().call(ctx, opUpdateAliases, &ChangeAliases{
		Actions: actions,
	}, resp)
	if err != nil {
		return newQdrantErr(err, "UpdateAliases")
	}
	return nil
}

// Creates a shard key for a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - collectionName: The name of the collection to create a shard key for.
//   - request: The CreateShardKey request containing shard key parameters.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) CreateShardKey(ctx context.Context, collectionName string, request *CreateShardKey) error {
	resp := &CreateShardKeyResponse{}
	err := c.get().call(ctx, opCreateShardKey, &CreateShardKeyRequest{
		CollectionName: collectionName,
		Request:        request,
	}, resp)
	if err != nil {
		return newQdrantErr(err, "CreateShardKey", collectionName)
	}
	return nil
}

// Deletes a shard key for a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - collectionName: The name of the collection to delete a shard key from.
//   - request: The DeleteShardKey request containing shard key parameters.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) DeleteShardKey(ctx context.Context, collectionName string, request *DeleteShardKey) error {
	resp := &DeleteShardKeyResponse{}
	err := c.get().call(ctx, opDeleteShardKey, &DeleteShardKeyRequest{
		CollectionName: collectionName,
		Request:        request,
	}, resp)
	if err != nil {
		return newQdrantErr(err, "DeleteShardKey", collectionName)
	}
	return nil
}

// Get cluster information for a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - collectionName: The name of the collection.
//
// Returns:
//   - *CollectionClusterInfoResponse: Detailed information about the collection's cluster.
//   - error: An error if the operation fails.
//
//nolint:lll	// Ignoring the long line length for naming consistency.
func (c *Client) GetCollectionClusterInfo(ctx context.Context, collectionName string) (*CollectionClusterInfoResponse, error) {
	response := &CollectionClusterInfoResponse{}
	err := c.get().call(ctx, opCollectionClusterInfo, &CollectionClusterInfoRequest{
		CollectionName: collectionName,
	}, response)
	if err != nil {
		return nil, newQdrantErr(err, "GetCollectionClusterInfo", collectionName)
	}
	return response, nil
}

// Update cluster setup for a collection.
//
// Parameters:
//   - ctx: The context for the request.
//   - request: The update to perform on the cluster setup.
//
// Returns:
//   - error: An error if the operation fails.
func (c *Client) UpdateClusterCollectionSetup(ctx context.Context, request *UpdateCollectionClusterSetupRequest) error {
	resp := &UpdateCollectionClusterSetupResponse{}
	err := c.get().call(ctx, opUpdateCollectionCluster, request, resp)
	if err != nil {
		return newQdrantErr(err, "UpdateClusterCollectionSetup", request.GetCollectionName())
	}
	return nil
}
