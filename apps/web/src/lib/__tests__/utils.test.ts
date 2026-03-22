import { describe, it, expect } from 'vitest';
import { formatLatency, formatBytes, cn } from '../utils';

describe('formatLatency', () => {
  it('formats sub-millisecond latency', () => {
    expect(formatLatency(0.5)).toBe('<1ms');
  });

  it('formats milliseconds', () => {
    expect(formatLatency(50)).toBe('50ms');
    expect(formatLatency(999)).toBe('999ms');
  });

  it('formats seconds', () => {
    expect(formatLatency(1000)).toBe('1.00s');
    expect(formatLatency(1500)).toBe('1.50s');
  });
});

describe('formatBytes', () => {
  it('formats zero bytes', () => {
    expect(formatBytes(0)).toBe('0 B');
  });

  it('formats bytes', () => {
    expect(formatBytes(512)).toBe('512 B');
  });

  it('formats kilobytes', () => {
    expect(formatBytes(1024)).toBe('1 KB');
    expect(formatBytes(1536)).toBe('1.5 KB');
  });

  it('formats megabytes', () => {
    expect(formatBytes(1048576)).toBe('1 MB');
  });
});

describe('cn', () => {
  it('merges class names', () => {
    expect(cn('foo', 'bar')).toBe('foo bar');
  });

  it('handles conditional classes', () => {
    expect(cn('foo', true && 'bar', false && 'baz')).toBe('foo bar');
  });

  it('handles tailwind classes', () => {
    expect(cn('text-red-500', 'text-blue-500')).toBe('text-blue-500');
  });
});
