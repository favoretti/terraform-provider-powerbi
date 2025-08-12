# Power BI Terraform Provider - Implementation Plan for Full API Coverage

## Current State Analysis

### Existing Resources
The provider currently implements the following resources:
1. **powerbi_workspace** - Manages Power BI workspaces (groups)
2. **powerbi_pbix** - Manages PBIX file uploads and updates
3. **powerbi_refresh_schedule** - Manages dataset refresh schedules
4. **powerbi_workspace_access** - Manages workspace user access
5. **powerbi_dataset** - Manages datasets

### Existing Data Sources
1. **powerbi_workspace** - Retrieves workspace information

### Current API Client Coverage
The provider has partial API implementations for:
- Groups (Workspaces) operations
- Datasets operations
- Reports operations
- Import operations
- Capacities operations
- Admin operations
- Push datasets operations
- Users operations

## Gap Analysis - Missing Resources

Based on the Power BI REST API v1.0 documentation, the following major resources are missing:

### 1. Dashboard Resources
**Priority: HIGH**
- `powerbi_dashboard` - Create and manage dashboards
- `powerbi_dashboard_tile` - Manage dashboard tiles
- Data sources for dashboards and tiles

### 2. App Resources
**Priority: MEDIUM**
- `powerbi_app` data source - Retrieve installed apps
- `powerbi_app_dashboard` data source - Access app dashboards
- `powerbi_app_report` data source - Access app reports

### 3. Gateway Resources
**Priority: HIGH**
- `powerbi_gateway` data source - Retrieve gateway information
- `powerbi_gateway_datasource` - Manage gateway data sources
- `powerbi_gateway_datasource_user` - Manage data source permissions

### 4. Dataflow Resources
**Priority: MEDIUM**
- `powerbi_dataflow` - Create and manage dataflows
- `powerbi_dataflow_refresh_schedule` - Manage dataflow refresh schedules
- Data sources for dataflows

### 5. Pipeline Resources
**Priority: MEDIUM**
- `powerbi_deployment_pipeline` - Manage deployment pipelines
- `powerbi_pipeline_stage` - Manage pipeline stages
- `powerbi_pipeline_operation` - Execute pipeline operations

### 6. Embed Token Resources
**Priority: LOW**
- `powerbi_embed_token` data source - Generate embed tokens for reports/dashboards

### 7. Template App Resources
**Priority: LOW**
- `powerbi_template_app` data source - Access template apps

## Implementation Roadmap

### Phase 1: Core Infrastructure (Weeks 1-2)
1. **API Client Enhancements**
   - Add missing API operations for dashboards
   - Add gateway API operations
   - Add dataflow API operations
   - Add pipeline API operations
   - Implement proper pagination for list operations
   - Add retry logic for rate limiting (429 responses)

### Phase 2: High Priority Resources (Weeks 3-6)
1. **Dashboard Management**
   ```go
   // New resources to implement:
   - resource_dashboard.go
   - resource_dashboard_tile.go
   - data_source_dashboard.go
   - data_source_dashboard_tiles.go
   ```

2. **Gateway Management**
   ```go
   // New resources to implement:
   - data_source_gateway.go
   - resource_gateway_datasource.go
   - resource_gateway_datasource_user.go
   ```

### Phase 3: Medium Priority Resources (Weeks 7-10)
1. **Dataflow Management**
   ```go
   // New resources to implement:
   - resource_dataflow.go
   - resource_dataflow_refresh_schedule.go
   - data_source_dataflow.go
   ```

2. **Deployment Pipeline Management**
   ```go
   // New resources to implement:
   - resource_deployment_pipeline.go
   - resource_pipeline_stage.go
   - resource_pipeline_operation.go
   ```

3. **App Data Sources**
   ```go
   // New data sources to implement:
   - data_source_app.go
   - data_source_app_dashboard.go
   - data_source_app_report.go
   ```

### Phase 4: Additional Features (Weeks 11-12)
1. **Enhanced Existing Resources**
   - Add missing properties to existing resources
   - Implement import functionality for all resources
   - Add validation for resource properties

2. **Embed Token Support**
   ```go
   // New data source to implement:
   - data_source_embed_token.go
   ```

3. **Template App Support**
   ```go
   // New data source to implement:
   - data_source_template_app.go
   ```

## Technical Implementation Details

### API Client Structure Improvements
```go
// Suggested new API client files:
internal/powerbiapi/
├── dashboards.go         // Dashboard operations
├── tiles.go              // Dashboard tile operations
├── gateways.go           // Gateway operations (enhance existing)
├── dataflows.go          // Dataflow operations
├── pipelines.go          // Deployment pipeline operations
├── apps.go               // App operations
├── embed.go              // Embed token operations
└── template_apps.go      // Template app operations
```

### Resource Schema Patterns
Each new resource should follow the existing patterns:
1. Implement CRUD operations (Create, Read, Update, Delete)
2. Support import functionality
3. Include proper error handling
4. Add comprehensive acceptance tests
5. Generate documentation using the existing docgen tool

### Testing Strategy
1. Unit tests for all new API client methods
2. Acceptance tests for each new resource
3. Integration tests for complex workflows
4. Mock API responses for reliable testing

## Priority Justification

**HIGH Priority:**
- Dashboards and Gateways are fundamental Power BI components
- Most enterprise deployments require these features
- High user demand based on GitHub issues

**MEDIUM Priority:**
- Dataflows and Pipelines are advanced features
- Used in mature Power BI deployments
- Apps are read-only operations

**LOW Priority:**
- Embed tokens are typically handled by applications
- Template apps have limited use cases

## Success Metrics

1. **Coverage Goal:** Achieve 90%+ coverage of Power BI REST API v1.0
2. **Quality Metrics:**
   - All resources have acceptance tests
   - Documentation generated for all resources
   - Import functionality for stateful resources
3. **Performance:**
   - Proper rate limiting handling
   - Efficient batch operations where applicable

## Breaking Changes Consideration

The implementation should maintain backward compatibility:
1. No changes to existing resource schemas
2. New optional fields only
3. Deprecation notices for any future removals
4. Version tagging for releases

## Maintenance and Updates

1. Monitor Power BI API changelog for new features
2. Quarterly review of API coverage
3. Community feedback integration
4. Regular dependency updates

## Estimated Timeline

- **Total Duration:** 12 weeks
- **Resources Required:** 1-2 developers
- **Testing/Documentation:** 20% of total effort

## Next Steps

1. Review and approve implementation plan
2. Set up development environment
3. Create GitHub issues for each phase
4. Begin Phase 1 implementation
5. Establish regular progress reviews