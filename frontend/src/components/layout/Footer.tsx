export async function Footer() {
  "use cache";
  return (
    <footer className="site-footer">
      <p className="footer-note">© {new Date().getFullYear()} Podlet. Built in the open.</p>
      <nav className="footer-nav" aria-label="Footer navigation">
        <a href="https://github.com/SudiMango/podlet#readme">Documentation</a>
        <a href="https://github.com/SudiMango/podlet">GitHub</a>
        <a href="https://github.com/SudiMango/podlet/issues">Feedback</a>
      </nav>
    </footer>
  );
}
