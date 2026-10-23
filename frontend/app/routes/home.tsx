import type { Route } from "../+types/root";
import { Welcome } from "../welcome/welcome";
import { Home } from "@/widgets/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "New React Router App" },
    { name: "description", content: "Welcome to React Router!" },
  ];
}

export default function HomePage() {
  return <Home />;
}
