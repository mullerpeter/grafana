import { DataFrameView } from '@grafana/data';
import { locationService } from '@grafana/runtime';

import { DashboardQueryResult, getGrafanaSearcher } from '../service';
import { SearchLayout } from '../types';
import * as utils from '../utils';

import { getSearchStateManager } from './SearchStateManager';

jest.mock('@grafana/runtime', () => {
  const originalModule = jest.requireActual('@grafana/runtime');
  return {
    ...originalModule,
  };
});

describe('SearchStateManager', () => {
  jest.spyOn(getGrafanaSearcher(), 'search').mockResolvedValue({
    isItemLoaded: jest.fn(),
    loadMoreItems: jest.fn(),
    totalRows: 0,
    view: new DataFrameView<DashboardQueryResult>({ fields: [], length: 0 }),
  });

  it('Can get search state manager with initial state', async () => {
    const stm = getSearchStateManager();
    expect(stm.state.layout).toBe(SearchLayout.Folders);
  });

  describe('initStateFromUrl', () => {
    it('should read and set state from URL and trigger search', async () => {
      const stm = getSearchStateManager();
      locationService.partial({ query: 'test', tag: ['tag1', 'tag2'] });
      stm.initStateFromUrl();
      expect(stm.state.folderUid).toBe(undefined);
      expect(stm.state.query).toBe('test');
      expect(stm.state.tag).toEqual(['tag1', 'tag2']);
    });

    it('should init or clear folderUid', async () => {
      const stm = getSearchStateManager();
      stm.initStateFromUrl('asdsadas');
      expect(stm.state.folderUid).toBe('asdsadas');

      stm.initStateFromUrl();
      expect(stm.state.folderUid).toBe(undefined);
    });

    it('should reset filters if state is updated and no URL params are present', () => {
      const parseRouteParamsSpy = jest.spyOn(utils, 'parseRouteParams');
      // Set initial values
      parseRouteParamsSpy.mockImplementation(() => ({
        query: 'hello',
        sort: 'alpha-asc',
      }));
      const stm = getSearchStateManager();
      stm.initStateFromUrl();

      // Verify that they have been set
      expect(stm.state.query).toBe('hello');
      expect(stm.state.sort).toBe('alpha-asc');
      expect(stm.state.folderUid).toBe(undefined);

      // Changed to a view with no URL state.
      parseRouteParamsSpy.mockImplementation(() => ({}));
      stm.initStateFromUrl('abc');

      expect(stm.state.query).toBe('');
      expect(stm.state.sort).toBe(undefined);
      expect(stm.state.folderUid).toBe('abc');
    });
  });
});
