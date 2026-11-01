import { redirect } from "react-router";
import { Initial } from "@/widgets/initial/ui";

export async function clientLoader() {
  if (window.sessionStorage.getItem("_tokens") !== null) {
    return redirect("/home");
  }

  return null;
}

export function HydrateFallback() {
  return <main className="min-h-screen bg-background" aria-busy="true" />;
}

export default function InitialPage() {
  return <Initial />;
}
