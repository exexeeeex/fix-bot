import { Button } from '@/shared/components/ui/button';
import { Icon } from '@/shared/components';
import { usePassQr } from '../../hooks';

export const Qr: React.FC = () => {
  const { qrCode, isLoading, isError, refetch } = usePassQr();

  return (
    <article className="w-full max-w-md overflow-hidden rounded-3xl border border-[#e5eaf2] bg-white shadow-[0_20px_60px_-24px_rgba(32,65,115,0.2)]">
      <header className="border-b border-[#e5eaf2] px-6 py-5 sm:px-8">
        <p className="text-sm font-semibold text-[#3478f6]">ФИКС</p>
        <h1 className="mt-1 text-2xl font-bold tracking-tight text-[#26344b]">Мой пропуск</h1>
        <p className="mt-2 text-sm leading-5 text-[#66758c]">
          Покажите QR-код на входе в учебное заведение, чтобы подтвердить посещение.
        </p>
      </header>

      <div className="flex min-h-80 flex-col items-center justify-center px-6 py-8">
        {isLoading ? (
          <div className="flex flex-col items-center gap-3 text-[#66758c]" role="status">
            <Icon
              icon="loader-circle"
              aria-hidden
              className="text-[#3478f6] [&_svg]:size-8 [&_svg]:animate-spin"
            />
            <p>Загружаем ваш QR-код...</p>
          </div>
        ) : isError ? (
          <div className="flex flex-col items-center text-center" role="alert">
            <Icon icon="alert-circle" aria-hidden className="text-red-600 [&_svg]:size-9" />
            <p className="mt-3 text-sm text-[#66758c]">
              Не удалось загрузить QR-код. Убедитесь, что вы вошли как студент, и попробуйте ещё
              раз.
            </p>
            <Button type="button" variant="outline" onClick={() => void refetch()} className="mt-5">
              <Icon icon="rotate-cw" aria-hidden className="mr-2 [&_svg]:size-4" />
              Повторить
            </Button>
          </div>
        ) : qrCode ? (
          <img
            src={qrCode}
            alt="QR-код электронного пропуска"
            className="size-73 max-w-full rounded-xl border border-[#e5eaf2] bg-white p-2"
          />
        ) : null}
      </div>
    </article>
  );
};
