import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import type { AuthenticationResponse } from '@/features/authentication/types';
import { useLogin } from '@/features/authentication/hooks';

type ApplicationStatus = AuthenticationResponse['user']['applicationStatus'];
type InitialStatus = Exclude<ApplicationStatus, 'approved'> | 'loading' | 'error';

export const useInitialForm = () => {
  const { login } = useLogin();
  const navigate = useNavigate();
  const [userId, setUserId] = useState('');
  const [name, setName] = useState<string>('');
  const [status, setStatus] = useState<InitialStatus>('loading');

  const authenticate = useCallback(async () => {
    setStatus('loading');

    try {
      const response = await login();
      if (!response.success) {
        throw new Error('Authentication was not successful');
      }

      setUserId(response.user.id);
      setName(`${response.user.first_name} ${response.user.last_name}`);
      if (response.user.applicationStatus === 'approved') {
        navigate('/home', { replace: true });
        return;
      }

      setStatus(response.user.applicationStatus);
    } catch (error: unknown) {
      console.error('Failed to check application status', error);
      setStatus('error');
    }
  }, [login, navigate]);

  useEffect(() => {
    void authenticate();
  }, [authenticate]);

  return {
    userId,
    name,
    status,
    retryAuthentication: authenticate,
  };
};
