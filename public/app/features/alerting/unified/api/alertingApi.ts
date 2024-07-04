import { BaseQueryFn, createApi } from '@reduxjs/toolkit/query/react';
import { lastValueFrom } from 'rxjs';

import { BackendSrvRequest, getBackendSrv } from '@grafana/runtime';

import { logMeasurement } from '../Analytics';

type RequestOptions = Omit<BackendSrvRequest, 'data'> & {
  /**
   * Data to send with a request. Maps to `data` property on a `BackendSrvRequest`
   *
   * This is done to allow us to more easily consume code-gen APIs that expect a `body` property
   */
  body?: BackendSrvRequest['data'];
};

export const backendSrvBaseQuery = (): BaseQueryFn<RequestOptions> => async (requestOptions) => {
  try {
    const requestStartTs = performance.now();

    const { body, ...options } = requestOptions;
    const modifiedOptions: BackendSrvRequest = {
      ...options,
      ...(body && { data: body }),
    };

    const { data, ...meta } = await lastValueFrom(getBackendSrv().fetch(modifiedOptions));

    logMeasurement(
      'backendSrvBaseQuery',
      {
        loadTimeMs: performance.now() - requestStartTs,
      },
      {
        url: requestOptions.url,
        method: requestOptions.method ?? 'GET',
        responseStatus: meta.statusText,
      }
    );

    return { data, meta };
  } catch (error) {
    return { error };
  }
};

export const alertingApi = createApi({
  reducerPath: 'alertingApi',
  baseQuery: backendSrvBaseQuery(),
  tagTypes: [
    'AlertingConfiguration',
    'AlertmanagerConfiguration',
    'AlertmanagerConnectionStatus',
    'AlertmanagerAlerts',
    'AlertmanagerSilences',
    'OnCallIntegrations',
    'OrgMigrationState',
    'DataSourceSettings',
    'GrafanaLabels',
    'CombinedAlertRule',
    'GrafanaRulerRule',
    'GrafanaSlo',
  ],
  endpoints: () => ({}),
});
