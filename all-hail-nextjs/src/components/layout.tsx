import React, { type ReactNode } from "react";
import Header from "./header";
import Sidebar from "./sidebar";

interface LayoutProps {
  children: ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
  return (
    <div>
      {/* <Header /> */}
      {/* <Sidebar /> */}
      <main>{children}</main>
    </div>
  );
};

export default Layout;
