import React from "react";

const Footer = () => {
  return (
    <footer className="pb-4shadow px-2 duration-500 ease-in-out md:px-5">
      <div className="mx-auto w-full max-w-screen-xl p-4 md:flex md:items-center md:justify-between">
        <span className="text-lg font-medium text-white dark:text-white sm:text-center">
          © 2023{" "}
          <a href="https://inviter.id/" className="hover:underline">
            Inviter™
          </a>
        </span>
        <ul className="mt-3 flex flex-wrap items-center text-lg font-medium text-white dark:text-white sm:mt-0">
          <li>
            <a href="#" className="mr-4 hover:underline md:mr-6 ">
              About
            </a>
          </li>
          <li>
            <a href="#" className="mr-4 hover:underline md:mr-6">
              Privacy Policy
            </a>
          </li>
          <li>
            <a href="#" className="mr-4 hover:underline md:mr-6">
              Licensing
            </a>
          </li>
          <li>
            <a href="#" className="hover:underline">
              Contact
            </a>
          </li>
        </ul>
      </div>
    </footer>
  );
};
export default Footer;
