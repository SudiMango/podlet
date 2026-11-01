import { Footer } from "@/components/layout/Footer";
import { Header } from "@/components/layout/Header";
import { Hero } from "@/components/sections/Hero";

export default function Home() {
  return (
    <div className="site-shell">
      <Header />
      <main>
        <Hero />
      </main>
      <Footer />
    </div>
  );
}
