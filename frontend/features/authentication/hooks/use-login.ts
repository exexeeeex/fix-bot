import { useCallback } from 'react';
import { getMax } from '@/shared/lib/platform/max';
import { storeAuthenticatedSession, useLoginMutation } from '../api/authentication';

export const useLogin = () => {
  const [triggerLogin, mutation] = useLoginMutation();

  const login = useCallback(async () => {
    const response = await triggerLogin({ initData: getMax().initData }).unwrap();
    storeAuthenticatedSession(response);
    return response;
  }, [triggerLogin]);

  return {
    login,
    isLoading: mutation.isLoading,
    isError: mutation.isError,
    error: mutation.error,
    isSuccess: mutation.isSuccess,
  };
};
