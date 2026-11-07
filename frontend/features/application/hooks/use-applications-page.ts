import { useEffect, useMemo, useState } from 'react';
import { useGetByOrganizationIdQuery, type Application } from '@/entities/application';
import type { Group } from '@/entities/group';
import { useGroupsByOrganizationQuery } from '@/entities/group';
import { useAppSelector } from '@/shared/lib/api/store/store';

const ALL_GROUPS = 'all';
const UNASSIGNED_GROUP = 'unassigned';

export interface GroupedApplications {
  id: string;
  name: string;
  applications: Application[];
}

export const useApplicationsPage = () => {
  const user = useAppSelector((state) => state.authentication.user);
  const [storedOrganizationId, setStoredOrganizationId] = useState('');
  const [groupId, setGroupId] = useState(ALL_GROUPS);
  const organizationId = user?.organization_id ?? storedOrganizationId;

  useEffect(() => {
    if (user?.organization_id) {
      setStoredOrganizationId(user.organization_id);
      return;
    }

    setStoredOrganizationId(window.sessionStorage.getItem('_organization-id') ?? '');
  }, [user?.organization_id]);

  const applicationsQuery = useGetByOrganizationIdQuery(organizationId, {
    skip: !organizationId,
  });
  const groupsQuery = useGroupsByOrganizationQuery(organizationId, { skip: !organizationId });
  const groupedApplications = useMemo(
    () => groupApplications(applicationsQuery.data ?? [], groupsQuery.data ?? []),
    [applicationsQuery.data, groupsQuery.data]
  );
  const groupOptions = useMemo(() => {
    if (groupsQuery.data?.length) {
      return groupsQuery.data.map((group) => ({ id: group.id, name: group.name }));
    }
    return groupedApplications.map(({ id, name }) => ({ id, name }));
  }, [groupsQuery.data, groupedApplications]);
  const visibleGroups =
    groupId === ALL_GROUPS
      ? groupedApplications
      : groupedApplications.filter((group) => group.id === groupId);
  const isLoading =
    (!!organizationId && (applicationsQuery.isLoading || groupsQuery.isLoading)) ||
    (!organizationId && !user);

  const pageState = !organizationId
    ? user
      ? 'no-organization'
      : 'loading-organization'
    : isLoading
      ? 'loading-applications'
      : applicationsQuery.isError
        ? 'applications-error'
        : 'content';

  return {
    organizationId,
    groupId,
    selectGroup: setGroupId,
    groupOptions,
    visibleGroups,
    pageState,
    isGroupsError: groupsQuery.isError,
    emptyMessage:
      groupId === ALL_GROUPS
        ? 'В этой организации пока нет заявок.'
        : 'В выбранной группе пока нет заявок.',
    retryApplications: applicationsQuery.refetch,
    retryGroups: groupsQuery.refetch,
  };
};

const groupApplications = (applications: Application[], groups: Group[]): GroupedApplications[] => {
  const groupNames = new Map(groups.map((group) => [group.id, group.name]));
  const grouped = new Map<string, Application[]>();

  for (const application of applications) {
    const groupId = application.group?.id ?? application.groupId;
    const key = groupId || UNASSIGNED_GROUP;
    const groupApplications = grouped.get(key) ?? [];
    groupApplications.push(application);
    grouped.set(key, groupApplications);
  }

  return Array.from(grouped, ([id, groupApplications]) => ({
    id,
    name:
      id === UNASSIGNED_GROUP
        ? 'Без группы'
        : (groupApplications[0].group?.name ?? groupNames.get(id) ?? `Группа ${id}`),
    applications: groupApplications,
  })).sort((first, second) => first.name.localeCompare(second.name, 'ru'));
};
