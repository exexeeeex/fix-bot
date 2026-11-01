import { type RouteConfig, index, layout, route } from "@react-router/dev/routes";

export default [
  index("routes/initial.tsx"),
  layout("routes/app-layout.tsx", [
    route("home", "routes/home.tsx"),
    route("profile", "routes/test.tsx"),
  ]),
] satisfies RouteConfig;
