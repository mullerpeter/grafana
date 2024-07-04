import { css } from '@emotion/css';
import { useMemo } from 'react';

import { GrafanaTheme2 } from '@grafana/data';
import { Button, LinkButton, LoadingPlaceholder, Stack, useStyles2 } from '@grafana/ui';
import {
  ALL_MUTE_TIMINGS,
  MuteTimingActionsButtons,
  useExportMuteTiming,
} from 'app/features/alerting/unified/components/mute-timings/MuteTimingActionsButtons';
import { MuteTimeInterval } from 'app/plugins/datasource/alertmanager/types';

import { Authorize } from '../../components/Authorize';
import { AlertmanagerAction, useAlertmanagerAbilities, useAlertmanagerAbility } from '../../hooks/useAbilities';
import { makeAMLink } from '../../utils/misc';
import { DynamicTable, DynamicTableColumnProps, DynamicTableItemProps } from '../DynamicTable';
import { EmptyAreaWithCTA } from '../EmptyAreaWithCTA';
import { ProvisioningBadge } from '../Provisioning';
import { Spacer } from '../Spacer';

import { useMuteTimings } from './useMuteTimings';
import { renderTimeIntervals } from './util';

interface MuteTimingsTableProps {
  alertManagerSourceName: string;
  hideActions?: boolean;
}

export const MuteTimingsTable = ({ alertManagerSourceName, hideActions }: MuteTimingsTableProps) => {
  const styles = useStyles2(getStyles);
  const [ExportAllDrawer, showExportAllDrawer] = useExportMuteTiming();

  const { data, isLoading } = useMuteTimings({ alertmanager: alertManagerSourceName });

  const items = useMemo((): Array<DynamicTableItemProps<MuteTimeInterval>> => {
    const muteTimings = data || [];

    return muteTimings.map((mute) => {
      return {
        id: mute.id,
        data: mute,
      };
    });
  }, [data]);

  const [_, allowedToCreateMuteTiming] = useAlertmanagerAbility(AlertmanagerAction.CreateMuteTiming);

  const [exportMuteTimingsSupported, exportMuteTimingsAllowed] = useAlertmanagerAbility(
    AlertmanagerAction.ExportMuteTimings
  );
  const columns = useColumns(alertManagerSourceName, hideActions);

  if (isLoading) {
    return <LoadingPlaceholder text="Loading mute timings..." />;
  }

  return (
    <div className={styles.container}>
      <Stack direction="row" alignItems="center">
        <span>
          Enter specific time intervals when not to send notifications or freeze notifications for recurring periods of
          time.
        </span>
        <Spacer />
        {!hideActions && items.length > 0 && (
          <Authorize actions={[AlertmanagerAction.CreateMuteTiming]}>
            <LinkButton
              className={styles.muteTimingsButtons}
              icon="plus"
              variant="primary"
              href={makeAMLink('alerting/routes/mute-timing/new', alertManagerSourceName)}
            >
              Add mute timing
            </LinkButton>
          </Authorize>
        )}
        {exportMuteTimingsSupported && (
          <>
            <Button
              icon="download-alt"
              className={styles.muteTimingsButtons}
              variant="secondary"
              aria-label="export all"
              disabled={!exportMuteTimingsAllowed}
              onClick={() => showExportAllDrawer(ALL_MUTE_TIMINGS)}
            >
              Export all
            </Button>
            {ExportAllDrawer}
          </>
        )}
      </Stack>
      {items.length > 0 ? (
        <DynamicTable items={items} cols={columns} pagination={{ itemsPerPage: 25 }} />
      ) : !hideActions ? (
        <EmptyAreaWithCTA
          text="You haven't created any mute timings yet"
          buttonLabel="Add mute timing"
          buttonIcon="plus"
          buttonSize="lg"
          href={makeAMLink('alerting/routes/mute-timing/new', alertManagerSourceName)}
          showButton={allowedToCreateMuteTiming}
        />
      ) : (
        <EmptyAreaWithCTA text="No mute timings configured" buttonLabel={''} showButton={false} />
      )}
    </div>
  );
};

function useColumns(alertManagerSourceName: string, hideActions = false) {
  const [[_editSupported, allowedToEdit], [_deleteSupported, allowedToDelete]] = useAlertmanagerAbilities([
    AlertmanagerAction.UpdateMuteTiming,
    AlertmanagerAction.DeleteMuteTiming,
  ]);
  const showActions = !hideActions && (allowedToEdit || allowedToDelete);

  return useMemo((): Array<DynamicTableColumnProps<MuteTimeInterval>> => {
    const columns: Array<DynamicTableColumnProps<MuteTimeInterval>> = [
      {
        id: 'name',
        label: 'Name',
        renderCell: function renderName({ data }) {
          return (
            <div>
              {data.name} {data.provisioned && <ProvisioningBadge tooltip />}
            </div>
          );
        },
        size: 1,
      },
      {
        id: 'timeRange',
        label: 'Time range',
        renderCell: ({ data }) => {
          return renderTimeIntervals(data);
        },
        size: 5,
      },
    ];
    if (showActions) {
      columns.push({
        id: 'actions',
        label: 'Actions',
        alignColumn: 'end',
        renderCell: ({ data }) => (
          <MuteTimingActionsButtons muteTiming={data} alertManagerSourceName={alertManagerSourceName} />
        ),
        size: 2,
      });
    }
    return columns;
  }, [showActions, alertManagerSourceName]);
}

const getStyles = (theme: GrafanaTheme2) => ({
  container: css({
    display: 'flex',
    flexFlow: 'column nowrap',
  }),
  muteTimingsButtons: css({
    marginBottom: theme.spacing(2),
  }),
});
