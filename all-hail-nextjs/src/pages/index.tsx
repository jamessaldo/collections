import type { ReactElement } from "react";
import Head from "next/head";
import Layout from "@/components/layout";
import NestedLayout from "@/components/nested-layout";
import { type NextPageWithLayout } from "./_app";
import Footer from "@/components/footer";
import Header from "@/components/header";

const Home: NextPageWithLayout = () => {
  return (
    <>
      <Head>
        <title>Inviter</title>
        <meta name="description" content="Making invitations easy and fun." />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <div className="bg-[conic-gradient(at_left,_var(--tw-gradient-stops))] from-yellow-300 via-pink-600 to-indigo-700">
        <Header />
        <main className="flex min-h-screen flex-col items-center justify-center">
          <div className="container flex flex-col items-center justify-center gap-12 px-4 py-16 ">
            <div className="max-w-screen-lg rounded-lg border-2 border-dashed border-gray-700 px-4 py-12 text-2xl">
              Inviter is a wedding invitation app that makes it easy and fun to
              create and send wedding invitations. With Inviter, you can choose
              from a variety of beautiful templates, customize your invitations
              with your own photos and text, and send them out to your guests
              with just a few clicks.
            </div>

            <h1 className="text-5xl tracking-tight text-white sm:text-[5rem]">
              Making invitations{" "}
              <span className="font-extrabold text-[#2B345F]">easy</span> and{" "}
              <span className="font-extrabold text-[#2B345F]">fun</span>.
            </h1>
            <div className="flex max-w-xl flex-col rounded-xl bg-white/10 p-4 text-white">
              <h3 className="text-2xl font-bold">Features</h3>
              <div className="text-lg">
                <ul className="list-inside list-disc">
                  <li>Choose from a variety of beautiful templates</li>
                  <li>
                    Customize your invitations with your own photos and text
                  </li>
                  <li>
                    Send invitations out to your guests with just a few clicks
                    (soon!)
                  </li>
                  <li>Track RSVPs</li>
                  <li>Share your invitations on social media</li>
                </ul>
              </div>
            </div>
            <div className="rounded-lg border-2 border-dashed border-gray-700 px-4 py-12 text-2xl">
              Start your wedding planning journey with Inviter today! Click the
              button below to get started.
            </div>
          </div>
        </main>
        <Footer />
      </div>
    </>
  );
};

Home.getLayout = function getLayout(page: ReactElement) {
  return (
    <Layout>
      <NestedLayout>{page}</NestedLayout>
    </Layout>
  );
};

export default Home;
