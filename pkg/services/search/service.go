package search

import (
	"context"
	"sort"

	"github.com/grafana/grafana/pkg/services/featuremgmt"

	"github.com/grafana/grafana/pkg/services/accesscontrol"

	"github.com/grafana/grafana/pkg/infra/db"
	"github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/dashboards/dashboardaccess"
	"github.com/grafana/grafana/pkg/services/search/model"
	"github.com/grafana/grafana/pkg/services/star"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
)

func ProvideService(cfg *setting.Cfg, sqlstore db.DB, starService star.Service, dashboardService dashboards.DashboardService, acService accesscontrol.Service, features featuremgmt.FeatureToggles) *SearchService {
	s := &SearchService{
		Cfg: cfg,
		sortOptions: map[string]model.SortOption{
			SortAlphaAsc.Name:  SortAlphaAsc,
			SortAlphaDesc.Name: SortAlphaDesc,
		},
		sqlstore:         sqlstore,
		starService:      starService,
		dashboardService: dashboardService,
		acService:        acService,
		features:         features,
	}
	return s
}

type Query struct {
	Title         string
	Tags          []string
	OrgId         int64
	SignedInUser  *user.SignedInUser
	Limit         int64
	Page          int64
	IsStarred     bool
	Type          string
	DashboardUIDs []string
	DashboardIds  []int64
	// Deprecated: use FolderUID instead
	FolderIds  []int64
	FolderUIDs []string
	Permission dashboardaccess.PermissionType
	Sort       string
}

type Service interface {
	SearchHandler(context.Context, *Query) (model.HitList, error)
	SortOptions() []model.SortOption
}

type SearchService struct {
	Cfg              *setting.Cfg
	sortOptions      map[string]model.SortOption
	sqlstore         db.DB
	starService      star.Service
	dashboardService dashboards.DashboardService
	acService        accesscontrol.Service
	features         featuremgmt.FeatureToggles
}

func (s *SearchService) SearchHandler(ctx context.Context, query *Query) (model.HitList, error) {
	starredQuery := star.GetUserStarsQuery{
		UserID: query.SignedInUser.UserID,
	}
	staredDashIDs, err := s.starService.GetByUser(ctx, &starredQuery)
	if err != nil {
		return nil, err
	}

	// No starred dashboards will be found
	if query.IsStarred && len(staredDashIDs.UserStars) == 0 {
		return model.HitList{}, nil
	}

	// filter by starred dashboard IDs when starred dashboards are requested and no UID or ID filters are specified to improve query performance
	if query.IsStarred && len(query.DashboardIds) == 0 && len(query.DashboardUIDs) == 0 {
		for id := range staredDashIDs.UserStars {
			query.DashboardIds = append(query.DashboardIds, id)
		}
	}

	userFolderUIDs := make([]string, 0)
	userDashboardsUIDs := make([]string, 0)
	if s.features.IsEnabled(ctx, featuremgmt.FlagRbacNG) {
		// Lookup user's folders and dashboards
		userFolderUIDs, err = s.acService.LookupResources(ctx, query.SignedInUser, accesscontrol.LookupQuery{
			ResourceType: accesscontrol.KindFolder,
			Permission:   accesscontrol.PermissionView,
		})
		if err != nil {
			return nil, err
		}
		userDashboardsUIDs, err = s.acService.LookupResources(ctx, query.SignedInUser, accesscontrol.LookupQuery{
			ResourceType: accesscontrol.KindDashboard,
			Permission:   accesscontrol.PermissionView,
		})
		if err != nil {
			return nil, err
		}
		// Pre-filter if user has access to small amount of dashboards
		if int64(len(userFolderUIDs)+len(userDashboardsUIDs)) < query.Limit {
			if len(query.DashboardUIDs) == 0 {
				query.DashboardUIDs = append(userDashboardsUIDs, userFolderUIDs...)
			}
		}
	}

	metrics.MFolderIDsServiceCount.WithLabelValues(metrics.Search).Inc()
	dashboardQuery := dashboards.FindPersistedDashboardsQuery{
		Title:         query.Title,
		SignedInUser:  query.SignedInUser,
		DashboardUIDs: query.DashboardUIDs,
		DashboardIds:  query.DashboardIds,
		Type:          query.Type,
		FolderIds:     query.FolderIds, // nolint:staticcheck
		FolderUIDs:    query.FolderUIDs,
		Tags:          query.Tags,
		Limit:         query.Limit,
		Page:          query.Page,
		Permission:    query.Permission,
	}

	if sortOpt, exists := s.sortOptions[query.Sort]; exists {
		dashboardQuery.Sort = sortOpt
	}

	hits, err := s.dashboardService.SearchDashboards(ctx, &dashboardQuery)
	if err != nil {
		return nil, err
	}

	if int64(len(userFolderUIDs)+len(userDashboardsUIDs)) >= query.Limit && s.features.IsEnabled(ctx, featuremgmt.FlagRbacNG) {
		// TODO: Do post-filtering if number of available dashboards is large
	}

	if query.Sort == "" {
		hits = sortedHits(hits)
	}

	// set starred dashboards
	for _, dashboard := range hits {
		if _, ok := staredDashIDs.UserStars[dashboard.ID]; ok {
			dashboard.IsStarred = true
		}
	}

	// filter for starred dashboards if requested
	if !query.IsStarred {
		return hits, nil
	}
	result := model.HitList{}
	for _, dashboard := range hits {
		if dashboard.IsStarred {
			result = append(result, dashboard)
		}
	}
	return result, nil
}

func sortedHits(unsorted model.HitList) model.HitList {
	hits := make(model.HitList, 0)
	hits = append(hits, unsorted...)

	sort.Sort(hits)

	for _, hit := range hits {
		sort.Strings(hit.Tags)
	}

	return hits
}
