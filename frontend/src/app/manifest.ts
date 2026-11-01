import type { MetadataRoute } from "next";

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "Podlet",
    short_name: "Podlet",
    description: "Deploy repositories to your own server from the command line.",
    start_url: "/",
    display: "standalone",
    background_color: "#b9c9d7",
    theme_color: "#25374c",
    icons: [
      { src: "/icon.svg", sizes: "any", type: "image/svg+xml", purpose: "any" },
      { src: "/apple-icon", sizes: "180x180", type: "image/png", purpose: "any" },
    ],
  };
}
