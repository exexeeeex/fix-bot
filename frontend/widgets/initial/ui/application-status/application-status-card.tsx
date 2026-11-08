import type { ReactNode } from 'react';

type ApplicationStatusCardProps = {
  title: string;
  description: string;
  icon: ReactNode;
  tone: 'info' | 'pending' | 'danger' | 'warning';
  action?: ReactNode;
};

const toneClasses = {
  info: {
    card: 'border-slate-200',
    header: 'bg-slate-800',
    icon: 'bg-slate-100 text-slate-700',
  },
  pending: {
    card: 'border-amber-200',
    header: 'bg-amber-700',
    icon: 'bg-amber-50 text-amber-700',
  },
  danger: {
    card: 'border-red-200',
    header: 'bg-red-700',
    icon: 'bg-red-50 text-red-700',
  },
  warning: {
    card: 'border-amber-200',
    header: 'bg-amber-700',
    icon: 'bg-amber-50 text-amber-700',
  },
};

export const ApplicationStatusCard = ({
  title,
  description,
  icon,
  tone,
  action,
}: ApplicationStatusCardProps) => {
  const colors = toneClasses[tone];

  return (
    <div
      className={`overflow-hidden rounded-3xl border bg-white text-left shadow-[0_20px_60px_-24px_rgba(32,65,115,0.22)] ${colors.card}`}
    >
      <div
        className={`relative overflow-hidden px-6 pb-7 pt-6 text-white sm:px-8 sm:pb-8 sm:pt-7 ${colors.header}`}
      >
        <div
          aria-hidden="true"
          className="pointer-events-none absolute -right-10 -top-14 size-40 rounded-full bg-white/10"
        />
        <p className="relative text-sm font-semibold tracking-wide text-white/75">ФИКС</p>
        <h1 className="relative mt-3 text-2xl font-bold tracking-tight sm:text-3xl">{title}</h1>
      </div>

      <div className="flex flex-col items-center px-6 py-8 text-center sm:px-8 sm:py-10">
        <div className={`flex size-16 items-center justify-center rounded-2xl ${colors.icon}`}>
          {icon}
        </div>
        <p className="mt-5 max-w-sm text-sm leading-6 text-[#66758c] sm:text-base">{description}</p>
        {action}
      </div>
    </div>
  );
};
