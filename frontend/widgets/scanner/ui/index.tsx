import { Button } from '@/shared/components/ui/button';
import { Icon } from '@/shared/components';
import { useScannedQr, useScannerCamera } from '../hooks';

export const Scanner = () => {
  const scannedQr = useScannedQr();
  const camera = useScannerCamera(scannedQr.read);
  const error = camera.error ?? scannedQr.error;

  const startScanning = () => {
    scannedQr.reset();
    void camera.start();
  };

  return (
    <main className="min-h-svh w-full px-4 py-4 sm:px-6 sm:py-6 lg:px-8">
      <section className="mx-auto flex min-h-[calc(100svh-6rem)] w-full flex-col mb-20 rounded-3xl border border-[#e5eaf2] bg-white p-5 shadow-[0_20px_60px_-24px_rgba(32,65,115,0.2)] sm:min-h-[calc(100svh-3rem)] sm:p-7 lg:min-h-[calc(100svh-3rem)] lg:p-8">
        <header className="flex items-center gap-3">
          <div className="flex size-11 shrink-0 items-center justify-center rounded-xl bg-blue-50 text-[#3478f6]">
            <Icon icon="qr" aria-hidden className="[&_svg]:size-6" />
          </div>
          <div>
            <h1 className="text-xl font-bold tracking-tight text-[#26344b] sm:text-2xl">
              Сканер пропусков
            </h1>
            <p className="mt-1 text-sm text-[#66758c]">
              Наведите камеру на QR-код пропуска. Следующий код будет считан автоматически.
            </p>
          </div>
        </header>

        <div className="mt-5 flex min-h-72 flex-1 overflow-hidden rounded-2xl bg-slate-950 sm:mt-6">
          <video
            ref={camera.videoRef}
            autoPlay
            muted
            playsInline
            onPlaying={camera.handleVideoPlaying}
            role="region"
            aria-label="Область сканирования QR-кода"
            className="min-h-72 w-full object-cover"
          />
        </div>

        <div className="mt-5 flex flex-wrap items-center gap-3">
          {camera.isScanning ? (
            <Button
              type="button"
              variant="outline"
              onClick={() => void camera.stop()}
              className="h-11 rounded-xl px-5"
            >
              Остановить камеру
            </Button>
          ) : (
            <Button
              type="button"
              onClick={startScanning}
              disabled={camera.isStarting || scannedQr.isProcessing}
              className="h-11 rounded-xl px-5"
            >
              {camera.isStarting || scannedQr.isProcessing ? (
                <>
                  <Icon
                    icon="loader-circle"
                    aria-hidden
                    className="mr-2 [&_svg]:size-4 [&_svg]:animate-spin"
                  />
                  {scannedQr.isProcessing ? 'Проверяем код...' : 'Запускаем камеру...'}
                </>
              ) : (
                'Начать сканирование'
              )}
            </Button>
          )}

          {scannedQr.result && (
            <span className="text-sm font-medium text-emerald-700" role="status">
              QR-код успешно распознан
            </span>
          )}
        </div>

        {error && (
          <div
            className="mt-4 flex items-start gap-2 rounded-xl bg-red-50 p-3 text-sm text-red-700"
            role="alert"
          >
            <Icon icon="alert-circle" aria-hidden className="mt-0.5 shrink-0 [&_svg]:size-4" />
            <p>{error}</p>
          </div>
        )}

        {scannedQr.result && (
          <div
            className="mt-4 rounded-2xl border border-emerald-200 bg-emerald-50 p-4 sm:p-5"
            aria-live="polite"
          >
            <h2 className="flex items-center gap-2 font-semibold text-emerald-800">
              <Icon icon="user" aria-hidden className="[&_svg]:size-5" />
              Данные пропуска
            </h2>
            <dl className="mt-3 grid gap-2 text-sm text-emerald-900 sm:grid-cols-2">
              <div>
                <dt className="font-medium">Группа</dt>
                <dd>{scannedQr.result.group.name}</dd>
              </div>
              <div>
                <dt className="font-medium">ID студента</dt>
                <dd className="break-all">{scannedQr.result.id}</dd>
              </div>
            </dl>
          </div>
        )}
      </section>
    </main>
  );
};
