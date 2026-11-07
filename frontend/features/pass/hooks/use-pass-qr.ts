import { useGetMyQrCodeQuery } from '@/entities/qr';

export const usePassQr = () => {
  const { data, isLoading, isError, refetch } = useGetMyQrCodeQuery();

  return {
    qrCode: data,
    isLoading,
    isError,
    refetch,
  };
};
