import type { Metadata, Viewport } from "next";
import localFont from "next/font/local";
import "./site.css";

const sans = localFont({
  src: "./fonts/UbuntuSans.ttf",
  variable: "--font-sans",
  display: "swap",
  weight: "100 900",
});

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_SITE_URL ?? "https://podlet.dev"),
  title: {
    default: "Podlet — Deploy to your own server",
    template: "%s — Podlet",
  },
  description: "Connect a repository and deploy to your own server from the command line with Podlet.",
  applicationName: "Podlet",
  authors: [{ name: "Podlet", url: "https://podlet.dev" }],
  creator: "Podlet",
  publisher: "Podlet",
  category: "developer tools",
  keywords: ["deployment", "self-hosting", "server", "developer tools", "CLI"],
  alternates: { canonical: "/" },
  manifest: "/manifest.webmanifest",
  openGraph: {
    type: "website",
    locale: "en_US",
    url: "/",
    siteName: "Podlet",
    title: "Podlet — Deploy to your own server",
    description: "Connect a repository and deploy to your own server from the command line.",
  },
  twitter: {
    card: "summary_large_image",
    title: "Podlet — Deploy to your own server",
    description: "Connect a repository and deploy to your own server from the command line.",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: { index: true, follow: true, "max-image-preview": "large", "max-snippet": -1 },
  },
  formatDetection: { email: false, address: false, telephone: false },
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  themeColor: "#25374c",
  colorScheme: "light",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" className={sans.variable}>
      <body>{children}</body>
    </html>
  );
}
