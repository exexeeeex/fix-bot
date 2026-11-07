import { useCallback } from 'react';
import { useUpdateStatusMutation } from '@/entities/application';

export const useAcceptApplication = () => {
  const [triggerUpdateStatus, mutation] = useUpdateStatusMutation();

  const accept = useCallback(
    (id: string) => triggerUpdateStatus({ id, status: 'approved' }),
    [triggerUpdateStatus]
  );

  return {
    accept,
    isLoading: mutation.isLoading,
    isError: mutation.isError,
  };
};
