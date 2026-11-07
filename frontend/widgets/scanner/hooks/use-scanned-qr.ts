import { useCallback, useRef, useState } from 'react';
import { useReadQrCodeMutation } from '@/entities/qr';
import type { ReadQrCodeResponse } from '@/entities/qr';

export const useScannedQr = () => {
  const [triggerRead] = useReadQrCodeMutation();
  const [result, setResult] = useState<ReadQrCodeResponse | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const pendingScansRef = useRef(new Set<string>());
  const requestIdRef = useRef(0);

  const read = useCallback(
    async (value: string) => {
      if (pendingScansRef.current.has(value)) return;
      pendingScansRef.current.add(value);
      const requestId = ++requestIdRef.current;
      setIsProcessing(true);
      setError(null);

      try {
        const response = await triggerRead({ value }).unwrap();
        if (requestId === requestIdRef.current) setResult(response);
      } catch (cause: unknown) {
        console.error('Failed to read scanned QR code', cause);
        if (requestId === requestIdRef.current) {
          setError('Не удалось проверить QR-код. Убедитесь, что это действующий пропуск.');
        }
      } finally {
        pendingScansRef.current.delete(value);
        if (pendingScansRef.current.size === 0) setIsProcessing(false);
      }
    },
    [triggerRead]
  );

  const reset = useCallback(() => {
    requestIdRef.current += 1;
    setResult(null);
    setError(null);
  }, []);

  return { read, reset, result, isProcessing, error };
};
