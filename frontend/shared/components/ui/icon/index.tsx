import { ICON_MAP, type IconType } from "@/shared/lib";

interface IconProps {
    icon: IconType;
    size?: number;
	className?: string;
} 

export const Icon: React.FC<IconProps> = ({icon, size = 20, className}) => {
    const spanIcon = ICON_MAP[icon]

    return (<span
			className={className}
			style={{
				fontSize: size,
				lineHeight: 1,
				display: "inline-block",
			}}
			role='img'
			aria-label={icon}
		>
			{spanIcon}
		</span>)
} 