import React, { ReactNode } from "react";

interface NestedLayoutProps {
  children: ReactNode;
}

const NestedLayout = ({ children }: NestedLayoutProps) => {
  return (
    <div>
      <div>{children}</div>
    </div>
  );
};

export default NestedLayout;
