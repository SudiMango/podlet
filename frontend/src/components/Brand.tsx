export function BrandMark({ className = "brand-mark" }: { className?: string }) {
  return (
    <svg className={className} width="32" height="32" viewBox="0 0 32 32" fill="none" aria-hidden="true">
      <path d="M12 3h8a9 9 0 0 1 9 9v8a9 9 0 0 1-9 9h-8a9 9 0 0 1-9-9v-8a9 9 0 0 1 9-9Z" stroke="currentColor" strokeWidth="2.5" />
      <path d="m9 11 5 5-5 5m9 0h5" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

export function Brand() {
  return (
    <span className="brand">
      <BrandMark />
      <span>podlet</span>
    </span>
  );
}
