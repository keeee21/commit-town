"use client";

import Image from "next/image";
import type { Session } from "next-auth";
import { signOutAction } from "@/libs/auth";

type DashboardPresenterProps = {
  session: Session;
};

export function DashboardPresenter({ session }: DashboardPresenterProps) {
  const { user } = session;

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 font-sans dark:bg-black">
      <main className="flex min-h-screen w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
        <div className="w-full flex justify-between items-center">
          <Image
            className="dark:invert"
            src="/next.svg"
            alt="Next.js logo"
            width={100}
            height={20}
            priority
          />
          <form action={signOutAction}>
            <button
              type="submit"
              className="rounded-md bg-zinc-900 px-4 py-2 text-sm text-white hover:bg-zinc-800 dark:bg-zinc-50 dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              Sign out
            </button>
          </form>
        </div>

        <div className="flex flex-col items-center gap-8 text-center sm:items-start sm:text-left w-full">
          {/* ユーザー情報カード */}
          <div className="w-full max-w-md rounded-lg border border-zinc-200 dark:border-zinc-800 p-6 bg-zinc-50 dark:bg-zinc-900">
            <div className="flex items-center gap-4 mb-4">
              {user.image && (
                <Image
                  src={user.image}
                  alt={user.name || "User"}
                  width={64}
                  height={64}
                  className="rounded-full"
                />
              )}
              <div className="flex-1">
                <h2 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                  {user.name || "Anonymous User"}
                </h2>
                {user.email && (
                  <p className="text-sm text-zinc-600 dark:text-zinc-400">
                    {user.email}
                  </p>
                )}
              </div>
            </div>

            <div className="space-y-2 text-sm">
              <div className="flex justify-between py-2 border-t border-zinc-200 dark:border-zinc-800">
                <span className="text-zinc-600 dark:text-zinc-400">
                  User ID:
                </span>
                <span className="font-mono text-zinc-900 dark:text-zinc-50">
                  {session.user.id || "N/A"}
                </span>
              </div>
              {user.name && (
                <div className="flex justify-between py-2 border-t border-zinc-200 dark:border-zinc-800">
                  <span className="text-zinc-600 dark:text-zinc-400">
                    Username:
                  </span>
                  <span className="font-mono text-zinc-900 dark:text-zinc-50">
                    {user.name}
                  </span>
                </div>
              )}
            </div>
          </div>

          <div className="w-full max-w-md">
            <h1 className="text-3xl font-semibold leading-10 tracking-tight text-black dark:text-zinc-50 mb-4">
              Welcome back, {user.name?.split(" ")[0] || "User"}!
            </h1>
            <p className="text-lg leading-8 text-zinc-600 dark:text-zinc-400">
              You are successfully authenticated with GitHub.
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}
