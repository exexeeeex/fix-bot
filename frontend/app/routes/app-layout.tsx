import { Navigation } from "@/widgets/navigation";
import { Outlet, redirect } from "react-router";

export async function clientLoader() {
  if (window.sessionStorage.getItem("_tokens") === null) {
    return redirect("/");
  }

  return null;
}

export function HydrateFallback() {
  return <main className="min-h-screen bg-background" aria-busy="true" />;
}

export default function AppLayout() {
  return (
    <>
      <Outlet />
      <Navigation />
    </>
  );
}
