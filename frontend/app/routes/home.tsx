import type { Route } from "../+types/root";
import { Home } from "@/widgets/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Сфера — пропуск" },
    { name: "description", content: "Электронный пропуск студента" },
  ];
}

export default function HomePage() {
  return <Home />;
}
