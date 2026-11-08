import { useEffect } from 'react';
import { useGetAbcensesQuery } from '../api';

export const useGetAbcenses = () => {
  let organizationId = sessionStorage.getItem('_organization-id') ?? '';

  const { data, isLoading, isError, error } = useGetAbcensesQuery(organizationId);

  return {
    data,
    isLoading,
    isError,
    error,
  };
};
