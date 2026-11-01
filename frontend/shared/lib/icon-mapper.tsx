import {
  ArrowRight,
  Book,
  CalendarDays,
  ChartColumn,
  GraduationCap,
  QrCode,
  SquarePlus,
  User,
} from 'lucide-react';

export type IconType =
  'user' | 'charts' | 'calendar' | 'qr' | 'book' | 'plus' | 'arrow-right' | 'graduation-cap';

export const ICON_MAP: Record<IconType, React.ReactNode> = {
  user: <User />,
  charts: <ChartColumn />,
  calendar: <CalendarDays />,
  qr: <QrCode />,
  book: <Book />,
  plus: <SquarePlus />,
  'arrow-right': <ArrowRight />,
  'graduation-cap': <GraduationCap />,
};
