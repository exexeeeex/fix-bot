import {
  ArrowRight,
  AlertCircle,
  Book,
  CalendarDays,
  ChartColumn,
  Clock3,
  GraduationCap,
  LoaderCircle,
  QrCode,
  RotateCw,
  SquarePlus,
  User,
} from 'lucide-react';

export type IconType =
  | 'user'
  | 'charts'
  | 'calendar'
  | 'qr'
  | 'book'
  | 'plus'
  | 'arrow-right'
  | 'graduation-cap'
  | 'alert-circle'
  | 'clock'
  | 'loader-circle'
  | 'rotate-cw';

export const ICON_MAP: Record<IconType, React.ReactNode> = {
  user: <User />,
  charts: <ChartColumn />,
  calendar: <CalendarDays />,
  qr: <QrCode />,
  book: <Book />,
  plus: <SquarePlus />,
  'arrow-right': <ArrowRight />,
  'graduation-cap': <GraduationCap />,
  'alert-circle': <AlertCircle />,
  clock: <Clock3 />,
  'loader-circle': <LoaderCircle />,
  'rotate-cw': <RotateCw />,
};
