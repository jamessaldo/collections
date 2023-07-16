import React from "react";

const Sidebar = () => {
  function setDark(val: string) {
    if (val === "dark") {
      document.documentElement.classList.add("dark");
      //   moon.classList.add("hidden");
      //   sun.classList.remove("hidden");
    } else {
      document.documentElement.classList.remove("dark");
      //   sun.classList.add("hidden");
      //   moon.classList.remove("hidden");
    }
  }
  return (
    <aside className="fixed z-50 flex h-screen w-64 -translate-x-48 transform bg-[#1c3144] transition duration-1000 ease-in-out ">
      {/* sidebar button  */}
      <div className="max-toolbar absolute -right-6 top-2 flex h-12 w-full translate-x-24 scale-x-0 transform items-center justify-between rounded-full border-4 border-white bg-[#1c3144]  transition duration-300 ease-in dark:border-[#0d1e2c]">
        <div className="flex items-center space-x-2 pl-4 ">
          <div>
            <div
              onClick={() => setDark("dark")}
              className="moon text-white hover:text-blue-500 dark:hover:text-[#38BDF8]"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={3}
                stroke="currentColor"
                className="h-6 w-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M21.752 15.002A9.718 9.718 0 0118 15.75c-5.385 0-9.75-4.365-9.75-9.75 0-1.33.266-2.597.748-3.752A9.753 9.753 0 003 11.25C3 16.635 7.365 21 12.75 21a9.753 9.753 0 009.002-5.998z"
                />
              </svg>
            </div>
            <div
              onClick={() => setDark("light")}
              //   onClick="setDark('light')"
              className="sun hidden text-white hover:text-blue-500 dark:hover:text-[#38BDF8]"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                className="h-6 w-6"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 3v2.25m6.364.386l-1.591 1.591M21 12h-2.25m-.386 6.364l-1.591-1.591M12 18.75V21m-4.773-4.227l-1.591 1.591M5.25 12H3m4.227-4.773L5.636 5.636M15.75 12a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0z"
                />
              </svg>
            </div>
          </div>
          <div className="text-white hover:text-blue-500 dark:hover:text-[#38BDF8]">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={3}
              stroke="currentColor"
              className="h-6 w-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0"
              />
            </svg>
          </div>
        </div>
        <div className="group flex items-center space-x-3 rounded-full bg-gradient-to-r from-indigo-500 via-purple-500 to-purple-500 py-1  pl-10 pr-2 text-white dark:from-cyan-500 dark:to-blue-500  ">
          <div className="mr-12 transform duration-300 ease-in-out">
            Inviter
          </div>
        </div>
      </div>
      <div
        onClick={() => setDark("dark")}
        // onClick="openNav()"
        className="absolute -right-6 top-2 flex transform rounded-full border-4 border-white bg-[#1c3144] p-3 text-white transition duration-500 ease-in-out hover:rotate-45 hover:bg-purple-500 dark:border-[#0d1e2c] dark:hover:bg-blue-500"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={3}
          stroke="currentColor"
          className="h-6 w-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"
          />
        </svg>
      </div>
      {/* MAX SIDEBAR */}
      <div className="max mt-20 hidden h-[calc(100vh)] w-full flex-col space-y-2 text-white">
        <div className="flex w-full transform flex-row items-center space-x-3 rounded-full bg-[#1c3144] p-2 pl-8 text-white duration-300 ease-in-out hover:ml-4 hover:text-purple-500 dark:hover:text-blue-500">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="h-6 w-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"
            />
          </svg>
          <div>Home</div>
        </div>
        <div className="flex w-full transform flex-row items-center space-x-3 rounded-full bg-[#1c3144] p-2 pl-8 text-white duration-300 ease-in-out hover:ml-4 hover:text-purple-500 dark:hover:text-blue-500">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            className="h-6 w-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5"
            />
          </svg>
          <div>Table</div>
        </div>
        <div className="flex w-full transform flex-row items-center space-x-3 rounded-full bg-[#1c3144] p-2 pl-8 text-white duration-300 ease-in-out hover:ml-4 hover:text-purple-500 dark:hover:text-blue-500">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            className="h-6 w-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M10.5 6a7.5 7.5 0 107.5 7.5h-7.5V6z"
            />
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M13.5 10.5H21A7.5 7.5 0 0013.5 3v7.5z"
            />
          </svg>
          <div>Graph</div>
        </div>
      </div>
      {/* MINI SIDEBAR */}
      <div className="mini mt-20 flex h-[calc(100vh)] w-full flex-col space-y-2">
        <div className="flex w-full transform justify-end rounded-full bg-[#1c3144] p-3 pr-5 text-white duration-300 ease-in-out hover:ml-4 hover:text-purple-500 dark:hover:text-blue-500">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="h-6 w-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"
            />
          </svg>
        </div>
        <div className="flex w-full transform justify-end rounded-full bg-[#1c3144] p-3 pr-5 text-white duration-300 ease-in-out hover:ml-4 hover:text-purple-500 dark:hover:text-blue-500">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            className="h-6 w-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M3.75 5.25h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5m-16.5 4.5h16.5"
            />
          </svg>
        </div>
        <div className="flex w-full transform justify-end rounded-full bg-[#1c3144] p-3 pr-5 text-white duration-300 ease-in-out hover:ml-4 hover:text-purple-500 dark:hover:text-blue-500">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            className="h-6 w-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M10.5 6a7.5 7.5 0 107.5 7.5h-7.5V6z"
            />
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M13.5 10.5H21A7.5 7.5 0 0013.5 3v7.5z"
            />
          </svg>
        </div>
      </div>
    </aside>
  );
};

export default Sidebar;
