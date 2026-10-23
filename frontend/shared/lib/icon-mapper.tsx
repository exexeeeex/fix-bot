import { Book, CalendarDays, ChartColumn, QrCode, SquarePlus, User } from "lucide-react"

export type IconType = 'user' | 'charts' | 'calendar' | 'qr' | 'book' | 'plus'

export const ICON_MAP: Record<IconType, React.ReactNode> = {
    user: <User />,
    charts: <ChartColumn/>,
    calendar: <CalendarDays/>,
    qr: <QrCode/>,
    book: <Book/>,
    plus: <SquarePlus/>
}