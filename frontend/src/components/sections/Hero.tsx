import { InstallCommand } from "@/components/ui/InstallCommand";
import { BrandMark } from "@/components/Brand";

const cliCommand = "curl -fsSL https://get.podlet.dev/cli/install.sh | bash";
const agentCommand = "curl -fsSL https://get.podlet.dev/agent/install.sh | sudo bash";

export function Hero() {
  return (
    <section className="hero" aria-labelledby="hero-title">
      <div className="hero-inner">
        <div className="hero-emblem"><BrandMark className="hero-mark" /></div>
        <h1 className="hero-title" id="hero-title">Your server.<br />Your deployments.</h1>
        <p className="hero-copy">
          From repository to running on your own server.<br className="desktop-break" />
          Deploy from your terminal. Keep control of your code.
        </p>
        <a className="hero-docs" href="https://github.com/SudiMango/podlet#readme"><span>Get started with Podlet</span><svg width="18" height="18" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" aria-hidden="true"><path d="M4 10h12m-5-5 5 5-5 5" /></svg></a>
        <div className="install-panel" id="install" aria-label="Install Podlet">
          <div className="install-panel-bar" aria-hidden="true">
            <span>Install Podlet</span>
            <span className="terminal-label">Terminal</span>
          </div>
          <InstallCommand label="Computer" command={cliCommand} />
          <InstallCommand label="Server" command={agentCommand} />
        </div>
        <p className="install-note">Currently supports Linux and macOS only, on ARM64 and x86_64.</p>
      </div>
    </section>
  );
}
