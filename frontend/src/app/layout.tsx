import type { Metadata } from "next";
import localFont from "next/font/local";
import "./site.css";

const sans = localFont({
  src: "./fonts/UbuntuSans.ttf",
  variable: "--font-sans",
  display: "swap",
  weight: "100 900",
});

export const metadata: Metadata = {
  title: "Podlet — Ship on your own server",
  description: "Connect a repository, deploy to your own server, and keep infrastructure simple with Podlet.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" className={sans.variable}>
      <body>{children}</body>
    </html>
  );
}
