import type { UserRole } from "@/entities/user";
import { LINKS } from "../config";
import { NavigationItem } from "./item";

export const Navigation = () => {

    const role: UserRole = 'user';

    const links = LINKS.filter((link) => {
        if (Array.isArray(link.access)) {
            return link.access.includes(role);
        }

        return link.access === role;
    });

    return (
        <nav className="fixed bottom-0 left-0 w-full h-15 bg-white dark:bg-gray-950  flex flex-col items-center justify-center">
            <ul className="flex flex-row justify-center items-center gap-5 w-full h-15">
                {links.map((link) => (
                    <NavigationItem item={link}/>
                ))}
            </ul>
        </nav>
    )
}
