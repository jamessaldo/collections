import React, { ReactNode } from "react";

interface NestedLayoutProps {
  children: ReactNode;
}

const NestedLayout = ({ children }: NestedLayoutProps) => {
  return (
    <div>
      <section>Sidebar</section>
      <div>{children}</div>
    </div>
  );
};

export default NestedLayout;
