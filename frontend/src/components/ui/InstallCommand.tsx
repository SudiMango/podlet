"use client";

import { useEffect, useRef, useState } from "react";

type InstallCommandProps = { label: string; command: string };

export function InstallCommand({ label, command }: InstallCommandProps) {
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState(false);
  const timeout = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => () => {
    if (timeout.current) clearTimeout(timeout.current);
  }, []);

  async function copyCommand() {
    try {
      await navigator.clipboard.writeText(command);
      setError(false);
      setCopied(true);
      if (timeout.current) clearTimeout(timeout.current);
      timeout.current = setTimeout(() => setCopied(false), 1600);
    } catch {
      setError(true);
    }
  }

  return (
    <div className="command-row">
      <span className="command-label">{label}</span>
      <code className="command-code" tabIndex={0} aria-label={`${label} install command, scroll horizontally to read`}><span className="prompt">$</span>{command}</code>
      <button
        className="copy-button"
        type="button"
        onClick={copyCommand}
        aria-label={copied ? `${label} command copied` : `Copy ${label} command`}
        title={copied ? "Copied" : "Copy command"}
      >
        {copied ? (
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true"><path d="m5 12 4 4L19 6" /></svg>
        ) : (
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" aria-hidden="true">
            <rect x="8" y="8" width="11" height="11" rx="2" />
            <path d="M16 8V6a2 2 0 0 0-2-2H6a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h2" />
          </svg>
        )}
      </button>
      <span className={error ? "copy-message" : "sr-only"} role="status">
        {error ? "Couldn't copy. Select the command and copy it manually." : copied ? `${label} command copied.` : ""}
      </span>
    </div>
  );
}
