import { useCallback, useEffect, useRef, useState } from 'react';
import { useZxing } from 'react-zxing';

const getErrorMessage = (error: unknown) =>
  error instanceof Error ? error.message : 'Неизвестная ошибка';

export const useScannerCamera = (onDecoded: (value: string) => void) => {
  const [isCameraActive, setIsCameraActive] = useState(false);
  const [isStarting, setIsStarting] = useState(false);
  const [isScanning, setIsScanning] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const onDecodedRef = useRef(onDecoded);
  const isCameraActiveRef = useRef(false);
  const isStartingRef = useRef(false);
  const lastDecodedValueRef = useRef<string | null>(null);

  useEffect(() => {
    onDecodedRef.current = onDecoded;
  }, [onDecoded]);

  const handleDecodeResult = useCallback((result: { rawValue: string }) => {
    if (!isCameraActiveRef.current || lastDecodedValueRef.current === result.rawValue) return;

    lastDecodedValueRef.current = result.rawValue;
    console.info('[QR scanner] QR-код распознан', {
      payloadLength: result.rawValue.length,
    });
    onDecodedRef.current(result.rawValue);
  }, []);

  const handleCameraError = useCallback((cause: unknown) => {
    console.error('Failed to start QR scanner camera', cause);
    isCameraActiveRef.current = false;
    isStartingRef.current = false;
    setIsCameraActive(false);
    setIsStarting(false);
    setIsScanning(false);
    setError(
      `Не удалось получить доступ к камере (${getErrorMessage(cause)}). Разрешите доступ к камере и попробуйте снова.`
    );
  }, []);

  const { ref: videoRef } = useZxing({
    paused: !isCameraActive,
    constraints: {
      audio: false,
      video: { facingMode: 'environment' },
    },
    formats: ['qr_code'],
    timeBetweenDecodingAttempts: 200,
    onDecodeResult: handleDecodeResult,
    onError: handleCameraError,
  });

  const handleVideoPlaying = useCallback(() => {
    if (!isCameraActiveRef.current) return;

    isStartingRef.current = false;
    setIsStarting(false);
    setIsScanning(true);
  }, []);

  const start = useCallback(() => {
    if (isStartingRef.current || isCameraActiveRef.current) return;

    lastDecodedValueRef.current = null;
    isStartingRef.current = true;
    isCameraActiveRef.current = true;
    setIsStarting(true);
    setError(null);
    setIsCameraActive(true);
  }, []);

  const stop = useCallback(() => {
    isCameraActiveRef.current = false;
    isStartingRef.current = false;
    lastDecodedValueRef.current = null;
    setIsCameraActive(false);
    setIsStarting(false);
    setIsScanning(false);
  }, []);

  return {
    videoRef,
    isStarting,
    isScanning,
    error,
    start,
    stop,
    handleVideoPlaying,
  };
};
