import { RegisterMyAbsenseButton } from '@/features/register-my-absence';
import { Qr } from '@/features/pass';

export const Home: React.FC = () => {
  return (
    <section className="mx-auto flex min-h-svh w-full max-w-3xl flex-col items-center gap-6 px-4 py-8 sm:justify-center sm:px-6">
      <Qr />
      <div className="flex justify-center">
        <RegisterMyAbsenseButton userId="1123" />
      </div>
    </section>
  );
};
