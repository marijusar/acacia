'use client';

import * as React from 'react';
import { cn } from '@/lib/utils';
import { useEscapeKey } from '@/hooks/use-escape-key';

interface SlidePanelProps {
  isOpen: boolean;
  onClose: () => void;
  side?: 'left' | 'right';
  width?: string;
  children: React.ReactNode;
}

export function SlidePanel({
  isOpen,
  onClose,
  width = '250px',
  children,
}: SlidePanelProps) {
  useEscapeKey(onClose, isOpen);

  if (!isOpen) return null;

  return (
    <div
      className={cn(
        'h-full sticky top-0 self-start bg-background border-l shadow-lg transition-all duration-300 ease-in-out flex-shrink-0 overflow-hidden',
        isOpen ? 'opacity-100' : 'opacity-0 pointer-events-none'
      )}
      style={{ width, maxHeight: 'calc(100vh - 64px)' }}
    >
      {children}
    </div>
  );
}
