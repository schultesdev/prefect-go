package client

import (
	"context"
	"fmt"
	"net/url"

	"git.schultes.dev/schultesdev/prefect-go/pkg/models"
	"github.com/google/uuid"
)

// BlockTypesService handles operations related to block types.
type BlockTypesService struct {
	client *Client
}

// Create creates a new block type.
func (s *BlockTypesService) Create(ctx context.Context, req *models.BlockTypeCreate) (*models.BlockType, error) {
	var blockType models.BlockType
	if err := s.client.post(ctx, "/block_types/", req, &blockType); err != nil {
		return nil, fmt.Errorf("failed to create block type: %w", err)
	}
	return &blockType, nil
}

// Get retrieves a block type by ID.
func (s *BlockTypesService) Get(ctx context.Context, id uuid.UUID) (*models.BlockType, error) {
	var blockType models.BlockType
	path := joinPath("/block_types", id.String())
	if err := s.client.get(ctx, path, &blockType); err != nil {
		return nil, fmt.Errorf("failed to get block type: %w", err)
	}
	return &blockType, nil
}

// GetBySlug retrieves a block type by slug.
func (s *BlockTypesService) GetBySlug(ctx context.Context, slug string) (*models.BlockType, error) {
	var blockType models.BlockType
	path := joinPath("/block_types/slug", slug)
	if err := s.client.get(ctx, path, &blockType); err != nil {
		return nil, fmt.Errorf("failed to get block type by slug: %w", err)
	}
	return &blockType, nil
}

// List retrieves block types with optional filtering.
func (s *BlockTypesService) List(ctx context.Context, filter *models.BlockTypeFilter, offset, limit int) ([]models.BlockType, error) {
	type nameFilter struct {
		Like *string `json:"like_,omitempty"`
	}
	type slugFilter struct {
		Any []string `json:"any_,omitempty"`
	}
	type blockTypeFilterCriteria struct {
		Name *nameFilter `json:"name,omitempty"`
		Slug *slugFilter `json:"slug,omitempty"`
	}
	type capabilitiesFilter struct {
		All []string `json:"all_,omitempty"`
	}
	type blockSchemaFilterCriteria struct {
		Capabilities *capabilitiesFilter `json:"block_capabilities,omitempty"`
	}
	type request struct {
		BlockTypes   *blockTypeFilterCriteria   `json:"block_types,omitempty"`
		BlockSchemas *blockSchemaFilterCriteria  `json:"block_schemas,omitempty"`
		Offset       int                        `json:"offset,omitempty"`
		Limit        int                        `json:"limit,omitempty"`
	}

	req := request{Offset: offset, Limit: limit}

	if filter != nil {
		btf := &blockTypeFilterCriteria{}
		if filter.Name != nil {
			btf.Name = &nameFilter{Like: filter.Name}
		}
		if len(filter.Slugs) > 0 {
			btf.Slug = &slugFilter{Any: filter.Slugs}
		}
		if btf.Name != nil || btf.Slug != nil {
			req.BlockTypes = btf
		}
		if len(filter.Capabilities) > 0 {
			req.BlockSchemas = &blockSchemaFilterCriteria{
				Capabilities: &capabilitiesFilter{All: filter.Capabilities},
			}
		}
	}

	var blockTypes []models.BlockType
	if err := s.client.post(ctx, "/block_types/filter", req, &blockTypes); err != nil {
		return nil, fmt.Errorf("failed to list block types: %w", err)
	}
	return blockTypes, nil
}

// Update updates a block type.
func (s *BlockTypesService) Update(ctx context.Context, id uuid.UUID, req *models.BlockTypeUpdate) error {
	path := joinPath("/block_types", id.String())
	if err := s.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to update block type: %w", err)
	}
	return nil
}

// Delete deletes a block type by ID.
func (s *BlockTypesService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/block_types", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete block type: %w", err)
	}
	return nil
}

// InstallSystemBlockTypes installs the system block types.
func (s *BlockTypesService) InstallSystemBlockTypes(ctx context.Context) error {
	if err := s.client.post(ctx, "/block_types/install_system_block_types", nil, nil); err != nil {
		return fmt.Errorf("failed to install system block types: %w", err)
	}
	return nil
}

// ListBlockDocumentsBySlug retrieves all block documents for a block type identified by slug.
func (s *BlockTypesService) ListBlockDocumentsBySlug(ctx context.Context, slug string, includeSecrets bool) ([]models.BlockDocument, error) {
	path := joinPath("/block_types/slug", slug, "block_documents")
	if includeSecrets {
		query := url.Values{}
		query.Set("include_secrets", "true")
		path = buildPathWithValues(path, query)
	}
	var blockDocuments []models.BlockDocument
	if err := s.client.get(ctx, path, &blockDocuments); err != nil {
		return nil, fmt.Errorf("failed to list block documents for block type: %w", err)
	}
	return blockDocuments, nil
}

// GetBlockDocumentByName retrieves a block document by name for a given block type slug.
func (s *BlockTypesService) GetBlockDocumentByName(ctx context.Context, slug, name string, includeSecrets bool) (*models.BlockDocument, error) {
	path := joinPath("/block_types/slug", slug, "block_documents/name", name)
	if includeSecrets {
		query := url.Values{}
		query.Set("include_secrets", "true")
		path = buildPathWithValues(path, query)
	}
	var blockDocument models.BlockDocument
	if err := s.client.get(ctx, path, &blockDocument); err != nil {
		return nil, fmt.Errorf("failed to get block document by name: %w", err)
	}
	return &blockDocument, nil
}

// BlockSchemasService handles operations related to block schemas.
type BlockSchemasService struct {
	client *Client
}

// Create creates a new block schema.
func (s *BlockSchemasService) Create(ctx context.Context, req *models.BlockSchemaCreate) (*models.BlockSchema, error) {
	var blockSchema models.BlockSchema
	if err := s.client.post(ctx, "/block_schemas/", req, &blockSchema); err != nil {
		return nil, fmt.Errorf("failed to create block schema: %w", err)
	}
	return &blockSchema, nil
}

// Get retrieves a block schema by ID.
func (s *BlockSchemasService) Get(ctx context.Context, id uuid.UUID) (*models.BlockSchema, error) {
	var blockSchema models.BlockSchema
	path := joinPath("/block_schemas", id.String())
	if err := s.client.get(ctx, path, &blockSchema); err != nil {
		return nil, fmt.Errorf("failed to get block schema: %w", err)
	}
	return &blockSchema, nil
}

// GetByChecksum retrieves a block schema by checksum and optional version.
func (s *BlockSchemasService) GetByChecksum(ctx context.Context, checksum string, version string) (*models.BlockSchema, error) {
	path := joinPath("/block_schemas/checksum", checksum)
	if version != "" {
		query := url.Values{}
		query.Set("version", version)
		path = buildPathWithValues(path, query)
	}
	var blockSchema models.BlockSchema
	if err := s.client.get(ctx, path, &blockSchema); err != nil {
		return nil, fmt.Errorf("failed to get block schema by checksum: %w", err)
	}
	return &blockSchema, nil
}

// List retrieves block schemas with optional filtering.
func (s *BlockSchemasService) List(ctx context.Context, filter *models.BlockSchemaFilter, offset, limit int) ([]models.BlockSchema, error) {
	type blockTypeIDFilter struct {
		Any []string `json:"any_,omitempty"`
	}
	type capabilitiesFilter struct {
		All []string `json:"all_,omitempty"`
	}
	type versionFilter struct {
		Any []string `json:"any_,omitempty"`
	}
	type blockSchemaFilterCriteria struct {
		BlockTypeID  *blockTypeIDFilter  `json:"block_type_id,omitempty"`
		Capabilities *capabilitiesFilter `json:"block_capabilities,omitempty"`
		Version      *versionFilter      `json:"version,omitempty"`
	}
	type request struct {
		BlockSchemas *blockSchemaFilterCriteria `json:"block_schemas,omitempty"`
		Offset       int                        `json:"offset,omitempty"`
		Limit        int                        `json:"limit,omitempty"`
	}

	req := request{Offset: offset, Limit: limit}

	if filter != nil {
		bsf := &blockSchemaFilterCriteria{}
		if filter.BlockTypeID != nil {
			bsf.BlockTypeID = &blockTypeIDFilter{Any: []string{filter.BlockTypeID.String()}}
		}
		if len(filter.Capabilities) > 0 {
			bsf.Capabilities = &capabilitiesFilter{All: filter.Capabilities}
		}
		if filter.Version != nil {
			bsf.Version = &versionFilter{Any: []string{*filter.Version}}
		}
		if bsf.BlockTypeID != nil || bsf.Capabilities != nil || bsf.Version != nil {
			req.BlockSchemas = bsf
		}
	}

	var blockSchemas []models.BlockSchema
	if err := s.client.post(ctx, "/block_schemas/filter", req, &blockSchemas); err != nil {
		return nil, fmt.Errorf("failed to list block schemas: %w", err)
	}
	return blockSchemas, nil
}

// Delete deletes a block schema by ID.
func (s *BlockSchemasService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/block_schemas", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete block schema: %w", err)
	}
	return nil
}

// BlockDocumentsService handles operations related to block documents.
type BlockDocumentsService struct {
	client *Client
}

// Create creates a new block document.
func (s *BlockDocumentsService) Create(ctx context.Context, req *models.BlockDocumentCreate) (*models.BlockDocument, error) {
	var blockDocument models.BlockDocument
	if err := s.client.post(ctx, "/block_documents/", req, &blockDocument); err != nil {
		return nil, fmt.Errorf("failed to create block document: %w", err)
	}
	return &blockDocument, nil
}

// Get retrieves a block document by ID.
func (s *BlockDocumentsService) Get(ctx context.Context, id uuid.UUID, includeSecrets bool) (*models.BlockDocument, error) {
	path := joinPath("/block_documents", id.String())
	if includeSecrets {
		query := url.Values{}
		query.Set("include_secrets", "true")
		path = buildPathWithValues(path, query)
	}
	var blockDocument models.BlockDocument
	if err := s.client.get(ctx, path, &blockDocument); err != nil {
		return nil, fmt.Errorf("failed to get block document: %w", err)
	}
	return &blockDocument, nil
}

// List retrieves block documents with optional filtering.
func (s *BlockDocumentsService) List(ctx context.Context, filter *models.BlockDocumentFilter, includeSecrets bool, sort models.BlockDocumentSort, offset, limit int) ([]models.BlockDocument, error) {
	type idFilter struct {
		Any []string `json:"any_,omitempty"`
	}
	type nameFilter struct {
		Any []string `json:"any_,omitempty"`
	}
	type boolFilter struct {
		Eq bool `json:"eq_"`
	}
	type blockDocumentFilterCriteria struct {
		BlockTypeID *idFilter   `json:"block_type_id,omitempty"`
		Name        *nameFilter `json:"name,omitempty"`
		IsAnonymous *boolFilter `json:"is_anonymous,omitempty"`
	}
	type request struct {
		BlockDocuments *blockDocumentFilterCriteria `json:"block_documents,omitempty"`
		IncludeSecrets bool                         `json:"include_secrets,omitempty"`
		Sort           models.BlockDocumentSort     `json:"sort,omitempty"`
		Offset         int                          `json:"offset,omitempty"`
		Limit          int                          `json:"limit,omitempty"`
	}

	req := request{
		IncludeSecrets: includeSecrets,
		Sort:           sort,
		Offset:         offset,
		Limit:          limit,
	}

	if filter != nil {
		bdf := &blockDocumentFilterCriteria{}
		if filter.BlockTypeID != nil {
			bdf.BlockTypeID = &idFilter{Any: []string{filter.BlockTypeID.String()}}
		}
		if filter.Name != nil {
			bdf.Name = &nameFilter{Any: []string{*filter.Name}}
		}
		if filter.IsAnonymous != nil {
			bdf.IsAnonymous = &boolFilter{Eq: *filter.IsAnonymous}
		}
		if bdf.BlockTypeID != nil || bdf.Name != nil || bdf.IsAnonymous != nil {
			req.BlockDocuments = bdf
		}
	}

	var blockDocuments []models.BlockDocument
	if err := s.client.post(ctx, "/block_documents/filter", req, &blockDocuments); err != nil {
		return nil, fmt.Errorf("failed to list block documents: %w", err)
	}
	return blockDocuments, nil
}

// Count returns the number of block documents matching the filter.
func (s *BlockDocumentsService) Count(ctx context.Context, filter *models.BlockDocumentFilter) (int, error) {
	type idFilter struct {
		Any []string `json:"any_,omitempty"`
	}
	type nameFilter struct {
		Any []string `json:"any_,omitempty"`
	}
	type boolFilter struct {
		Eq bool `json:"eq_"`
	}
	type blockDocumentFilterCriteria struct {
		BlockTypeID *idFilter   `json:"block_type_id,omitempty"`
		Name        *nameFilter `json:"name,omitempty"`
		IsAnonymous *boolFilter `json:"is_anonymous,omitempty"`
	}
	type request struct {
		BlockDocuments *blockDocumentFilterCriteria `json:"block_documents,omitempty"`
	}

	var req request
	if filter != nil {
		bdf := &blockDocumentFilterCriteria{}
		if filter.BlockTypeID != nil {
			bdf.BlockTypeID = &idFilter{Any: []string{filter.BlockTypeID.String()}}
		}
		if filter.Name != nil {
			bdf.Name = &nameFilter{Any: []string{*filter.Name}}
		}
		if filter.IsAnonymous != nil {
			bdf.IsAnonymous = &boolFilter{Eq: *filter.IsAnonymous}
		}
		if bdf.BlockTypeID != nil || bdf.Name != nil || bdf.IsAnonymous != nil {
			req.BlockDocuments = bdf
		}
	}

	var count int
	if err := s.client.post(ctx, "/block_documents/count", req, &count); err != nil {
		return 0, fmt.Errorf("failed to count block documents: %w", err)
	}
	return count, nil
}

// Update updates a block document's data.
func (s *BlockDocumentsService) Update(ctx context.Context, id uuid.UUID, req *models.BlockDocumentUpdate) error {
	path := joinPath("/block_documents", id.String())
	if err := s.client.patch(ctx, path, req, nil); err != nil {
		return fmt.Errorf("failed to update block document: %w", err)
	}
	return nil
}

// Delete deletes a block document by ID.
func (s *BlockDocumentsService) Delete(ctx context.Context, id uuid.UUID) error {
	path := joinPath("/block_documents", id.String())
	if err := s.client.delete(ctx, path); err != nil {
		return fmt.Errorf("failed to delete block document: %w", err)
	}
	return nil
}

// BlockCapabilitiesService handles operations related to block capabilities.
type BlockCapabilitiesService struct {
	client *Client
}

// List retrieves all available block capabilities.
func (s *BlockCapabilitiesService) List(ctx context.Context) ([]string, error) {
	var capabilities []string
	if err := s.client.get(ctx, "/block_capabilities/", &capabilities); err != nil {
		return nil, fmt.Errorf("failed to list block capabilities: %w", err)
	}
	return capabilities, nil
}
