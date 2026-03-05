import { cn } from '../lib/utils';

interface StatusBadgeProps {
  state: 'up' | 'down' | 'degraded' | 'unknown';
  size?: 'sm' | 'md' | 'lg';
  showLabel?: boolean;
}

const stateConfig = {
  up: { 
    color: 'bg-emerald-500/80', 
    glow: 'shadow-emerald-500/30',
    label: 'Online',
  },
  down: { 
    color: 'bg-red-500/80', 
    glow: 'shadow-red-500/30',
    label: 'Offline',
  },
  degraded: { 
    color: 'bg-amber-500/80', 
    glow: 'shadow-amber-500/30',
    label: 'Degraded',
  },
  unknown: { 
    color: 'bg-neutral-500/80', 
    glow: 'shadow-neutral-500/30',
    label: 'Unknown',
  },
};

export function StatusBadge({ state, size = 'md', showLabel = false }: StatusBadgeProps) {
  const config = stateConfig[state];
  
  const sizeClasses = {
    sm: 'w-2.5 h-2.5',
    md: 'w-3 h-3',
    lg: 'w-3.5 h-3.5',
  };

  return (
    <div className="flex items-center gap-2.5">
      <div className="relative">
        {/* Outer glow ring - very subtle */}
        <div className={cn(
          'absolute inset-0 rounded-full animate-pulse-subtle',
          config.color,
          config.glow,
          'shadow-lg'
        )} />
        
        {/* Main dot */}
        <div className={cn(
          'relative rounded-full shadow-lg overflow-hidden',
          sizeClasses[size],
          config.color,
          'transition-all duration-500'
        )}>
          {/* Highlight overlay */}
          <div className="absolute inset-0 bg-gradient-to-br from-white/30 to-transparent" />
        </div>
      </div>
      
      {showLabel && (
        <span className="text-sm text-white/40">{config.label}</span>
      )}
    </div>
  );
}

export function getStatusColor(state: string): string {
  switch (state) {
    case 'up':
      return 'text-emerald-400/80';
    case 'down':
      return 'text-red-400/80';
    case 'degraded':
      return 'text-amber-400/80';
    default:
      return 'text-neutral-400/80';
  }
}

export function getStateBorder(state: string): string {
  switch (state) {
    case 'up':
      return 'border-emerald-500/15 hover:border-emerald-400/30';
    case 'down':
      return 'border-red-500/15 hover:border-red-400/30';
    case 'degraded':
      return 'border-amber-500/15 hover:border-amber-400/30';
    default:
      return 'border-neutral-500/15 hover:border-neutral-400/30';
  }
}
