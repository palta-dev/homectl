interface HeaderProps {
  onSearch?: (query: string) => void;
}

export function Header({ onSearch }: HeaderProps) {
  return (
    <header className="glass-header">
      <div className="container mx-auto px-6 py-4">
        <div className="flex items-center justify-between gap-4">
          {/* Logo and Title */}
          <div className="flex items-center gap-4">
            <div className="relative group">
              <div className="absolute inset-0 bg-white/10 rounded-2xl blur-xl opacity-20 group-hover:opacity-30 transition-opacity" />
              <div className="relative flex items-center justify-center w-12 h-12 rounded-2xl bg-gradient-to-br from-white/10 to-white/3 backdrop-blur-xl border border-white/8 shadow-2xl">
                <span className="text-lg font-semibold text-white/90">H</span>
              </div>
            </div>
          </div>

          {/* Search */}
          {onSearch && (
            <div className="hidden md:block relative">
              <input
                type="search"
                placeholder="Search services..."
                className="glass-input h-11 w-80 pl-10 pr-10 text-sm text-white/90 placeholder:text-white/30"
                onChange={(e) => onSearch(e.target.value)}
                aria-label="Search services"
              />
              <svg
                className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-white/30"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={1.5}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
          )}
        </div>
      </div>
    </header>
  );
}
