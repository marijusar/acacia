import type { ReactNode } from 'react';

type Props = React.HTMLAttributes<HTMLHeadingElement>;

export const Heading4 = ({ children, className }: Props) => {
  const css = `scroll-m-20 text-md font-semibold tracking-tight ${className}`;
  return <h4 className={css}>{children}</h4>;
};

export const Heading1 = ({ children, className }: Props) => {
  const css = `scroll-m-20 text-xl font-semibold tracking-tight ${className}`;
  return <h4 className={css}>{children}</h4>;
};
