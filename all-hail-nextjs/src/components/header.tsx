import React from "react";
import { signIn, signOut, useSession } from "next-auth/react";
import { api } from "@/utils/api";

const Header = () => {
  return (
    <header>
      <div className="mx-auto w-full max-w-screen-xl p-4 px-16 md:flex md:items-center md:justify-between">
        <span className="text-4xl font-extrabold dark:text-[#2B345F] sm:text-center">
          <a href="https://inviter.id/" className="hover:underline">
            Inviter
          </a>
        </span>
        <div className="dark:text-[#2B345F] sm:text-center">
          <AuthShowcase />
        </div>
      </div>
    </header>
  );
};

function AuthShowcase() {
  const { data: sessionData } = useSession();
  const message = api.greeting.hello.useQuery(
    {
      text: sessionData?.user?.name ?? "friend!",
    },
    { enabled: sessionData?.user !== undefined, staleTime: Infinity }
  );

  return (
    <div className="flex items-center justify-center gap-4">
      {sessionData && sessionData.user ? (
        <span className="text-xl font-semibold dark:text-[#2B345F]">
          {message.data ? message.data : "Loading tRPC query..."}
        </span>
      ) : (
        ""
      )}
      <button
        className="rounded-full bg-[#2B345F] px-10 py-3 font-semibold text-white no-underline transition hover:bg-[#2B345F]/80"
        onClick={sessionData ? () => void signOut() : () => void signIn()}
      >
        {sessionData ? "Sign out" : "Sign in"}
      </button>
    </div>
  );
}

export default Header;
