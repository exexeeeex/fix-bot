import { useCallback } from 'react';
import { useCreateMutation, type CreateApplicationRequest } from '@/entities/application';

export const useCreateApplication = () => {
  const [triggerCreate, mutation] = useCreateMutation();

  const create = useCallback(
    (data: CreateApplicationRequest) => triggerCreate(data),
    [triggerCreate]
  );

  return {
    create,
    isLoading: mutation.isLoading,
    isError: mutation.isError,
    isSuccess: mutation.isSuccess,
  };
};
