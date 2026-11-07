import { ICON_MAP, type IconType } from '@/shared/lib';

interface IconProps {
  icon: IconType;
  size?: number;
  className?: string;
  'aria-hidden'?: boolean;
}

export const Icon: React.FC<IconProps> = ({ icon, size, className, 'aria-hidden': ariaHidden }) => {
  const spanIcon = ICON_MAP[icon];

  return (
    <span
      className={className}
      style={{
        fontSize: size,
        lineHeight: 1,
        display: 'inline-block',
      }}
      role={ariaHidden ? undefined : 'img'}
      aria-label={ariaHidden ? undefined : icon}
      aria-hidden={ariaHidden}
    >
      {spanIcon}
    </span>
  );
};
