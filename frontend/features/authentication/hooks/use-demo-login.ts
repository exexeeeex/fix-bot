import { useCallback } from 'react';
import { storeAuthenticatedSession, useDemoLoginMutation } from '../api/authentication';

export const useDemoLogin = () => {
  const [triggerDemoLogin, mutation] = useDemoLoginMutation();

  const login = useCallback(
    async (role: 'student' | 'teacher') => {
      const response = await triggerDemoLogin({ role }).unwrap();
      storeAuthenticatedSession(response);
      return response;
    },
    [triggerDemoLogin]
  );

  return {
    login,
    isLoading: mutation.isLoading,
    isError: mutation.isError,
  };
};
