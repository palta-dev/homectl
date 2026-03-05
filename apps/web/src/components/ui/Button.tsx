import { forwardRef, ButtonHTMLAttributes } from 'react';
import { cn } from '../../lib/utils';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'ghost' | 'outline' | 'glass';
  size?: 'sm' | 'md' | 'lg' | 'icon';
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'glass', size = 'md', ...props }, ref) => {
    return (
      <button
        className={cn(
          'inline-flex items-center justify-center rounded-xl font-medium',
          'transition-all duration-300',
          'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
          'disabled:pointer-events-none disabled:opacity-50',
          {
            'bg-primary text-primary-foreground hover:bg-primary/90 shadow-lg shadow-primary/30': variant === 'default',
            'hover:bg-white/10 hover:text-accent-foreground': variant === 'ghost',
            'border border-white/20 bg-white/5 hover:bg-white/10 backdrop-blur-sm': variant === 'outline',
            'glass-button': variant === 'glass',
          },
          {
            'h-9 px-3 text-xs': size === 'sm',
            'h-11 px-4 py-2 text-sm': size === 'md',
            'h-12 px-6 text-base': size === 'lg',
            'h-11 w-11': size === 'icon',
          },
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);

Button.displayName = 'Button';
