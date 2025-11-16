import NextAuth from "next-auth";
import GitHub from "next-auth/providers/github";
import { envConfig } from "./env.config";

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    GitHub({
      clientId: envConfig.auth.github.clientId,
      clientSecret: envConfig.auth.github.clientSecret,
    }),
  ],
  callbacks: {
    async session({ session, token }) {
      if (token?.sub) {
        session.user.id = token.sub;
      }
      return session;
    },
  },
});
