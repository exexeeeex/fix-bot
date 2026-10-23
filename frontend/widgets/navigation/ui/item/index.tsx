import { NavLink } from "react-router";
import type { Link } from "../../config/_links";
import { Icon } from "@/shared/components";

interface NavigationItemProps {
    item: Link;
}

export const NavigationItem = ({ item }: NavigationItemProps) => {
    return (
        <li className="flex items-center justify-center">
            <NavLink to={item.href} end>
                {({ isActive }) => (
                    <div
                        className={`flex flex-col items-center justify-center ${
                            isActive ? "text-primary" : "text-muted"
                        }`}
                    >
                        <Icon
                            icon={item.icon}
                            className={ isActive ? "text-primary bg-primary/20 px-3 py-0.5 rounded-2xl" : "text-muted"}
                        />
                        <span>{item.name}</span>
                    </div>
                )}
            </NavLink>
        </li>
    );
};